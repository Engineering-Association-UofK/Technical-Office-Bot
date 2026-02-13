package database

import (
	"log/slog"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMySQLConnection(dsn string) (*sqlx.DB, error) {
	slog.Info("Starting MySQL connection...")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		slog.Error("Failed to open database connection: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	slog.Info("Opening initial SQL file...")
	file, err := os.ReadFile("migrations/000001_init.up.sql")
	if err != nil {
		return nil, err
	}
	slog.Info("Running initial SQL script...")

	queries := strings.Split(string(file), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	slog.Info("MySQL connection ready.")
	return db, nil
}
