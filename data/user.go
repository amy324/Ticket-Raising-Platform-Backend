package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	//"time"
)

// SetDB sets the database connection for the package
func SetDB(database *sql.DB) {
	db = database
}

// GetDB returns the current database connection
func GetDB() *sql.DB {
	return db
}

// InitializeDB initializes the database connection
func InitializeDB(user, password, host, port, dbname string) error {
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	var err error
	db, err = sql.Open("mysql", connectionStr)
	if err != nil {
		return err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// User structure
type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
}

func (u *User) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `
	INSERT INTO users (email, first_name, last_name, password, user_active, is_admin, pin_number)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password, // Store plaintext password
		u.Active,
		u.IsAdmin,
		"0000", // Replace with your desired default value
	)

	if err != nil {
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	newID = int(lastInsertID)

	return newID, nil
}

// GetByEmail retrieves a user by email
func GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, first_name, last_name, password, user_active, is_admin
		FROM users
		WHERE email = ?`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.IsAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// PasswordMatches compares the provided password with the stored password
func (u *User) PasswordMatches(providedPassword string) (bool, error) {
	return u.Password == providedPassword, nil
}

// AuthenticateUser authenticates a user based on email and password
func AuthenticateUser(email, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	// Check if the provided password matches the stored password
	matches, err := user.PasswordMatches(password)
	if err != nil {
		return nil, fmt.Errorf("error comparing passwords: %w", err)
	}

	if !matches {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
