package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"python-runner/configuration"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Global connection instance
var GlobalConnection *Connection

// Connection represents a MySQL database connection
type Connection struct {
	DB     *sqlx.DB
	config configuration.MySQLConfig
}

// ConnectionInterface defines the interface for database connections
type ConnectionInterface interface {
	Ping() error
	Close() error
	GetDB() *sqlx.DB
	IsConnected() bool
	Reconnect() error
}

// NewConnection creates a new MySQL database connection
func NewConnection() (*Connection, error) {
	config := configuration.GetMySQLConfig()

	conn := &Connection{
		config: config,
	}

	if err := conn.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	return conn, nil
}

// connect establishes the database connection
func (c *Connection) connect() error {
	dsn := configuration.GetMySQLConnectionString()

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	c.DB = db

	// Test the connection
	if err := c.Ping(); err != nil {
		c.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

// Ping tests the database connection
func (c *Connection) Ping() error {
	if c.DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return c.DB.Ping()
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB == nil {
		return nil
	}

	log.Printf("Closing MySQL database connection")
	return c.DB.Close()
}

// GetDB returns the underlying sqlx.DB instance
func (c *Connection) GetDB() *sqlx.DB {
	return c.DB
}

// IsConnected checks if the database connection is active
func (c *Connection) IsConnected() bool {
	if c.DB == nil {
		return false
	}

	if err := c.Ping(); err != nil {
		log.Printf("Database connection check failed: %v", err)
		return false
	}

	return true
}

// Reconnect attempts to reconnect to the database
func (c *Connection) Reconnect() error {
	log.Printf("Attempting to reconnect to MySQL database")

	// Close existing connection if it exists
	if c.DB != nil {
		c.Close()
	}

	return c.connect()
}

// GetStats returns database connection statistics
func (c *Connection) GetStats() sql.DBStats {
	if c.DB == nil {
		return sql.DBStats{}
	}
	return c.DB.Stats()
}

// Exec executes a query without returning any rows
func (c *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return c.DB.Exec(query, args...)
}

// Query executes a query that returns rows
func (c *Connection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return c.DB.Query(query, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (c *Connection) QueryRow(query string, args ...interface{}) *sql.Row {
	if c.DB == nil {
		// Return a row that will have an error when scanned
		return &sql.Row{}
	}
	return c.DB.QueryRow(query, args...)
}

// InitializeGlobalConnection initializes the global database connection
func InitializeGlobalConnection() error {
	conn, err := NewConnection()
	if err != nil {
		return fmt.Errorf("failed to initialize global MySQL connection: %w", err)
	}

	GlobalConnection = conn
	return nil
}

// GetGlobalConnection returns the global database connection
func GetGlobalConnection() *Connection {
	return GlobalConnection
}

// CloseGlobalConnection closes the global database connection
func CloseGlobalConnection() error {
	if GlobalConnection != nil {
		return GlobalConnection.Close()
	}
	return nil
}
