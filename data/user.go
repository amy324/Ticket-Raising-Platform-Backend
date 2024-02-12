package data

import (
	"context"

	"database/sql"
	"errors"
	"fmt"
	"log"

	"crypto/rand"
	"math/big"

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

var ErrUserNotFound = errors.New("user not found")


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

	// Generate a random 6-digit pin number for verification
	pinNumber, err := GeneratePinNumber()
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `
    INSERT INTO users (email, first_name, last_name, password, pin_number, user_active, is_admin, refreshJWT)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		pinNumber,
		u.UserActive,
		u.IsAdmin,
		refreshJWT,
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

func GeneratePinNumber() (string, error) {
	// Generate a random 6-digit number
	num, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", num), nil
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

func CreateAccessToken(userID int, userEmail string, accessToken string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Calculate the expiration time (30 minutes from creation)
	expirationTime := time.Now().Add(30 * time.Minute)

	// Check if the access token already exists for the user
	var existingAccessToken string
	err := db.QueryRowContext(ctx, "SELECT accessJWT FROM access_tokens WHERE user_id = ?", userID).Scan(&existingAccessToken)

	if err == sql.ErrNoRows {
		// If no rows are found, insert the access token for the user with expiration time
		result, err := db.ExecContext(ctx, "INSERT INTO access_tokens (user_id, email, accessJWT, created_at, expires_at) VALUES (?, ?, ?, NOW(), ?)", userID, userEmail, accessToken, expirationTime)
		if err != nil {
			return 0, err
		}

		return result.LastInsertId()
	} else if err != nil {
		return 0, err
	}

	// If the access token already exists, update it
	_, err = db.ExecContext(ctx, "UPDATE access_tokens SET accessJWT = ?, created_at = NOW(), expires_at = ? WHERE user_id = ?", accessToken, expirationTime, userID)
	if err != nil {
		return 0, err
	}

	// Return a dummy LastInsertId, as it's not relevant for updates
	return 0, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT id, email, first_name, last_name, password, user_active, is_admin
        FROM users
        WHERE id = ?`

	var user User
	row := db.QueryRowContext(ctx, query, userID)

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

// UpdateRefreshToken updates or inserts the refresh token for a user ID
func UpdateRefreshToken(userID int, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Check if the refresh token already exists for the user
	var existingRefreshToken string
	err := db.QueryRowContext(ctx, "SELECT refreshJWT FROM users WHERE id = ?", userID).Scan(&existingRefreshToken)

	if err == sql.ErrNoRows {
		// If no rows are found, insert the refresh token for the user
		_, err := db.ExecContext(ctx, "UPDATE users SET refreshJWT = ? WHERE id = ?", refreshToken, userID)
		return err
	} else if err != nil {
		return err
	}

	// If the refresh token already exists, update it
	_, err = db.ExecContext(ctx, "UPDATE users SET refreshJWT = ? WHERE id = ?", refreshToken, userID)
	return err
}

func (u *User) Logout() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Delete the access token entry in the access_tokens table
	_, err := db.ExecContext(ctx, "DELETE FROM access_tokens WHERE user_id = ?", u.ID)
	if err != nil {
		return err
	}

	// Set the refreshJWT to an empty string in the users table
	_, err = db.ExecContext(ctx, "UPDATE users SET refreshJWT = '' WHERE id = ?", u.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserIDByAccessToken retrieves the user ID associated with the given access token
func GetUserIDByAccessToken(accessToken string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var userID int
	query := `SELECT user_id FROM access_tokens WHERE accessJWT = ?`

	// Log the query being executed
	log.Printf("Executing query to retrieve user ID for access token: %s", accessToken)

	// Execute the query and scan the result
	err := db.QueryRowContext(ctx, query, accessToken).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return 0 if no rows are found
			return 0, nil
		}
		return 0, err
	}

	// Log the retrieved user ID
	log.Printf("Retrieved user ID from database: %d", userID)

	return userID, nil
}

// GetPinByEmail retrieves the PIN for a user by their email
func GetPinByEmail(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var pin string
	query := `
        SELECT pin_number FROM users WHERE email = ?`

	err := db.QueryRowContext(ctx, query, email).Scan(&pin)
	if err != nil {
		return "", err
	}

	return pin, nil
}

// UserExists checks if a user already exists in the database by email
func UserExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var count int
	query := `
        SELECT COUNT(*) FROM users WHERE email = ?`

	err := db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ActivateAccount activates the user account by setting UserActive to 1
func (u *User) ActivateAccount() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Update the user's isActive status to 1
	_, err := db.ExecContext(ctx, "UPDATE users SET user_active = 1 WHERE id = ?", u.ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePinAfterVerification updates the pin_number field after PIN verification
// and activates the user account
func (u *User) UpdatePinAfterVerification() error {
	// Prepare the SQL statement to update the pin_number and user_active fields
	query := "UPDATE users SET pin_number = ?, user_active = 1 WHERE id = ?"

	// Execute the SQL statement
	result, err := db.Exec(query, "N/A - verified", u.ID)
	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}

	// Check the number of rows affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected, user not found or pin update failed")
	}

	return nil
}

// GetAccessTokenExpirationTime retrieves the expiration time of the access token for the given user ID.
func GetAccessTokenExpirationTime(userID int) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var expiresAt time.Time
	err := db.QueryRowContext(ctx, "SELECT expires_at FROM access_tokens WHERE user_id = ?", userID).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a zero time if no access token found for the user
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return expiresAt.UTC(), nil // Ensure expiration time is in UTC
}

// GetUserEmailByAccessToken retrieves the email associated with the provided access token.
func GetUserEmailByAccessToken(accessToken string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var userEmail string
	err := db.QueryRowContext(ctx, "SELECT email FROM access_tokens WHERE accessJWT = ?", accessToken).Scan(&userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return an empty string if no email found for the access token
			return "", nil
		}
		return "", err
	}

	return userEmail, nil
}
// GetUserEmailByID retrieves the email associated with the provided user ID.
func GetUserEmailByID(userID int) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    var userEmail string
    err := db.QueryRowContext(ctx, "SELECT email FROM users WHERE id = ?", userID).Scan(&userEmail)
    if err != nil {
        if err == sql.ErrNoRows {
            // Return an empty string if no email found for the user ID
            return "", nil
        }
        return "", err
    }

    return userEmail, nil
}
