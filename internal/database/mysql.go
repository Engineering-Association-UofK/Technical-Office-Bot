package database

import (
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMySQLConnection(dsn string) *sqlx.DB {
	log.Println("Starting MySQL connection...")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Panic("Failed to open database connection: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		log.Panic("Database not responding: ", err)
	}

	log.Println("Opening initial SQL file...")
	file, err := os.ReadFile("migrations/000001_init.up.sql")
	if err != nil {
		log.Panic("Error reading initial sql file: ", err)
	}
	log.Println("Running starting SQL script...")

	queries := strings.Split(string(file), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err = db.Exec(query)
		if err != nil {
			log.Panic("Failed to run initial SQL script: ", err)
		}
	}

	log.Println("MySQL connection ready.")
	return db
}
