package main

import (
	"context"
	"database/sql"
	"fin-stat/api"
	"fin-stat/migrations"
	"fin-stat/transaction"
	"log"
	"net/http"
	"os"
	"time"
)

const commandTimeout = 5 * time.Second

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
	_, err = db.ExecContext(migrationCtx, migrations.InitialMigration)
	cancel()
	if err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	repo := transaction.NewPostgresRepository(db)
	svc := transaction.NewService(repo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           api.NewHandler(svc),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       commandTimeout,
		WriteTimeout:      commandTimeout,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("REST API listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("start HTTP server: %v", err)
	}
}
