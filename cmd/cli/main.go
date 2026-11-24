package main

import (
	"context"
	"database/sql"
	"fin-stat/cli"
	"fin-stat/migrations"
	"fin-stat/transaction"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

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

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(pingCtx)
	cancel()
	if err != nil {
		log.Fatalf("ping db: %v", err)
	}

	migrationCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = db.ExecContext(migrationCtx, migrations.InitialMigration)
	cancel()
	if err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	repo := transaction.NewPostgresRepository(db)
	svc := transaction.NewService(repo)

	cli.Run(svc)
}
