package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"gochat/models"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a shared database instance
var DB *sql.DB

// Connect initializes the database connection
func Connect(dbPath string) error {
	var err error

	// Open SQLite database
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Create tables if they don't exist
	if err = createTables(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// createTables creates required tables if they don't exist
func createTables() error {
	// Create users table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		status TEXT DEFAULT 'offline',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	return err
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
	mu sync.RWMutex // for thread safety
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Prepare statement
	stmt, err := r.db.Prepare(`
		INSERT INTO users (username, email, password, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare insert user statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	result, err := stmt.Exec(user.Username, user.Email, user.Password, user.Status, now, now)
	if err != nil {
		return fmt.Errorf("execute insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var user models.User
	err := r.db.QueryRow(`
		SELECT id, username, email, password, status, created_at, updated_at
		FROM users
		WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var user models.User
	err := r.db.QueryRow(`
		SELECT id, username, email, password, status, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("query user by id: %w", err)
	}

	return &user, nil
}

// UpdateUserStatus updates a user's status
func (r *UserRepository) UpdateUserStatus(id int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`
		UPDATE users
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, status, id)

	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}

	return nil
}
