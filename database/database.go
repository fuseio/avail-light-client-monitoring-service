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

// Add this struct after the Database struct
type ClientInfo struct {
	Address   string    `json:"address"`
	TokenID   string    `json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
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
		address TEXT PRIMARY KEY,
		token_id TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) RegisterClient(address, tokenID string) error {
	_, err := d.db.Exec(`
		INSERT INTO clients (address, token_id) 
		VALUES (?, ?)
		ON CONFLICT(address) DO UPDATE SET token_id = excluded.token_id
	`, address, tokenID)

	if err != nil {
		d.logger.Printf("Error registering client: %v", err)
		return err
	}
	return nil
}

func (d *Database) ClientExists(address string) (bool, error) {
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clients WHERE address = ?)", address).Scan(&exists)
	if err != nil {
		d.logger.Printf("Error checking client existence: %v", err)
		return false, err
	}
	return exists, nil
}

func (d *Database) GetClientTokenID(address string) (string, error) {
	var tokenID string
	err := d.db.QueryRow("SELECT token_id FROM clients WHERE address = ?", address).Scan(&tokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		d.logger.Printf("Error getting client token ID: %v", err)
		return "", err
	}
	return tokenID, nil
}

// Add this new function at the end of the file
func (d *Database) GetAllClients() ([]ClientInfo, error) {
	rows, err := d.db.Query(`
		SELECT address, token_id, created_at 
		FROM clients 
		ORDER BY created_at DESC
	`)
	if err != nil {
		d.logger.Printf("Error querying clients: %v", err)
		return nil, err
	}
	defer rows.Close()

	var clients []ClientInfo
	for rows.Next() {
		var client ClientInfo
		err := rows.Scan(&client.Address, &client.TokenID, &client.CreatedAt)
		if err != nil {
			d.logger.Printf("Error scanning client row: %v", err)
			return nil, err
		}
		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		d.logger.Printf("Error iterating client rows: %v", err)
		return nil, err
	}

	return clients, nil
}
