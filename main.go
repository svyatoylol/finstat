package main

import (
	"bufio"
	"context"
	"database/sql"
	_ "embed"
	"fin-stat/commands"
	"fin-stat/transaction"
	"fmt"
	"log"
	"os"
	"time"
)

const commandTimeout = 5 * time.Second

//go:embed migrations/001_init.sql
var initialMigration string

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://finstat:finstat@localhost:5432/finstat?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// sql.Open соединение не устанавливает, поэтому проверяем его явно
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(pingCtx)
	cancel()
	if err != nil {
		log.Fatalf("ping db: %v", err)
	}

	migrationCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = db.ExecContext(migrationCtx, initialMigration)
	cancel()
	if err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	repo := transaction.NewPostgresRepository(db)
	svc := transaction.NewService(repo)
	sess := commands.NewSession(0)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands (type 'exit' to quit):")

	for scanner.Scan() {
		cmd := commands.ParseCommand(scanner.Text())
		if cmd.Name == "exit" {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
		commands.ExecuteCommand(ctx, cmd, svc, sess)
		cancel() // без defer: мы внутри цикла

		fmt.Println("Enter a command:")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read input: %v", err)
	}
}
