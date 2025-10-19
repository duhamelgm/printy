package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection and operations
type Database struct {
	db *sql.DB
}

// New creates a new database connection and initializes the schema
func New(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	database := &Database{db: db}

	// Initialize schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	log.Printf("✅ Database initialized at: %s", dbPath)
	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// initSchema creates the database tables
func (d *Database) initSchema() error {
	// Create tickets table
	ticketsSQL := `
	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ref_id TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		priority INTEGER NOT NULL DEFAULT 0,
		cooldown INTEGER NOT NULL DEFAULT 0,
		weekdays TEXT NOT NULL DEFAULT '',
		assignee TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := d.db.Exec(ticketsSQL); err != nil {
		return fmt.Errorf("failed to create tickets table: %v", err)
	}

	// Add weekdays column if it doesn't exist (migration)
	alterSQL := `ALTER TABLE tickets ADD COLUMN weekdays TEXT DEFAULT '';`
	d.db.Exec(alterSQL) // Ignore error if column already exists

	// Add assignee column if it doesn't exist (migration)
	alterAssigneeSQL := `ALTER TABLE tickets ADD COLUMN assignee TEXT DEFAULT '';`
	d.db.Exec(alterAssigneeSQL) // Ignore error if column already exists

	// Create prints table
	printsSQL := `
	CREATE TABLE IF NOT EXISTS prints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ticket_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE
	);`

	if _, err := d.db.Exec(printsSQL); err != nil {
		return fmt.Errorf("failed to create prints table: %v", err)
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_tickets_ref_id ON tickets(ref_id);",
		"CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);",
		"CREATE INDEX IF NOT EXISTS idx_prints_ticket_id ON prints(ticket_id);",
		"CREATE INDEX IF NOT EXISTS idx_prints_created_at ON prints(created_at);",
	}

	for _, indexSQL := range indexes {
		if _, err := d.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
	}

	return nil
}

// GetDB returns the underlying database connection (for advanced operations)
func (d *Database) GetDB() *sql.DB {
	return d.db
}
