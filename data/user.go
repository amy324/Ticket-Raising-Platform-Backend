package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	//"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
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
	ID         int
	Email      string
	FirstName  string
	LastName   string
	Password   string
	PinNumber  string
	UserActive int
	IsAdmin    int
	RefreshJWT string
}

// AccessToken structure
type AccessToken struct {
	ID        int
	UserID    int
	Email     string
	AccessJWT string
}

func (u *User) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}

	// Generate a refresh token
	refreshJWT, err := generateRefreshToken(u.ID)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `
    INSERT INTO users (email, first_name, last_name, password, pin_number, user_active, is_admin, refresh_jwt)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		u.PinNumber,
		u.UserActive,
		u.IsAdmin,
		refreshJWT, // Ensure this is the last parameter
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

// Helper function to generate a refresh token
func generateRefreshToken(userID int) (string, error) {
	// Set the expiration time for the token (you can customize this)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(userID), // Include user ID in the refresh token claims
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key (replace with your own secret key)
	refreshJWT := []byte(os.Getenv("JWT_REFRESH_KEY"))
	if len(refreshJWT) == 0 {
		log.Fatal("JWT_REFRESH_KEY is not set in the environment")
	}
	signedToken, err := token.SignedString(refreshJWT)
	if err != nil {
		return "", err
	}

	return signedToken, nil
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
		&user.UserActive,
		&user.IsAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// PasswordMatches compares the provided password with the stored hash
func (u *User) PasswordMatches(providedPassword string) (bool, error) {
	hashedPassword := []byte(u.Password)
	providedPasswordBytes := []byte(providedPassword)

	err := bcrypt.CompareHashAndPassword(hashedPassword, providedPasswordBytes)
	if err == nil {
		return true, nil
	} else if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	} else {
		return false, err
	}
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

// CreateAccessToken function
func CreateAccessToken(userID int, email string, accessJWT string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
        INSERT INTO access_tokens (user_id, email, accessJWT)
        VALUES (?, ?, ?)`

	res, err := db.ExecContext(ctx, stmt, userID, email, accessJWT)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	newID := int(lastInsertID)

	return newID, nil
}
