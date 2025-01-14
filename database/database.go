package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db     *sql.DB
	logger *log.Logger
}

func NewDatabase(dbPath string, logger *log.Logger) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite only supports one writer
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize tables
	if err := initializeTables(db); err != nil {
		return nil, err
	}

	return &Database{
		db:     db,
		logger: logger,
	}, nil
}

func initializeTables(db *sql.DB) error {
	// Add your table creation statements here
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT NOT NULL,
		status TEXT NOT NULL,
		last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) RegisterClient(address string) error {
	query := `
		INSERT INTO clients (address, status) 
		VALUES (?, 'active')
	`
	_, err := d.db.Exec(query, address)
	return err
}

func (d *Database) ClientExists(address string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM clients WHERE address = ?)`
	err := d.db.QueryRow(query, address).Scan(&exists)
	return exists, err
}
