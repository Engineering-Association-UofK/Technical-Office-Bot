package database

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMySQLConnection() (*sqlx.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.App.DBUser,
		config.App.DBPassword,
		config.App.DBHost,
		config.App.DBPort,
		config.App.DBName,
	)

	slog.Info("Starting MySQL connection...")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %s", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database: %s", err)
	}

	schema, err := os.ReadFile("resources/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("Failed to read schema file: %s", err)
	}
	slog.Info("Running initial SQL script...")
	err = executeQuerry(string(schema), db)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute schema script: %s", err)
	}

	data, err := os.ReadFile("resources/data.sql")
	if err != nil {
		return nil, fmt.Errorf("Failed to read data file: %s", err)
	}
	slog.Info("Running data SQL script...")
	err = executeQuerry(string(data), db)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute data script: %s", err)
	}

	slog.Info("MySQL connection ready.")
	return db, nil
}

func executeQuerry(schema string, db *sqlx.DB) error {
	queries := strings.Split(string(schema), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
