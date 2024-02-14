// token_handlers.go

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var db *sql.DB
var dbTimeout = 5 * time.Second

// RefreshTokenHandler handles the refreshing of access tokens
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the refresh token from the Authorization header
	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Remove the "Bearer " prefix from the token
	refreshToken = refreshToken[len("Bearer "):]

	// Validate the refresh token and get the user information
	user, err := validateRefreshJWT(refreshToken, os.Getenv("JWT_REFRESH_KEY"))
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	accessToken, err := generateAuthJWT(user, os.Getenv("JWT_ACCESS_KEY"), 30*time.Minute)
	if err != nil {
		fmt.Println("Error generating access JWT token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update the access token and its expiration time in the database
	err = updateAccessToken(db, dbTimeout, user.ID, accessToken, 30*time.Minute)
	if err != nil {
		fmt.Println("Error updating access token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	response := map[string]interface{}{"message": "Token refreshed successfully", "accessToken": accessToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// updateAccessToken updates the access token and its expiration time in the database
func updateAccessToken(db *sql.DB, timeout time.Duration, userID int, newAccessToken string, validityDuration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Calculate the new expiration time
	expiresAt := time.Now().Add(validityDuration)

	stmt := `
        UPDATE access_tokens
        SET accessJWT = ?, expires_at = ?
        WHERE user_id = ?`

	_, err := db.ExecContext(ctx, stmt, newAccessToken, expiresAt, userID)
	return err
}

// refreshAccessToken is a helper function for refreshing the access token
func refreshAccessToken(w http.ResponseWriter, r *http.Request, refreshToken string, db *sql.DB) {
	// Validate the refresh token and get the user information
	user, err := validateRefreshJWT(refreshToken, os.Getenv("JWT_REFRESH_KEY"))
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	accessToken, err := generateAuthJWT(user, os.Getenv("JWT_ACCESS_KEY"), 30*time.Minute)
	if err != nil {
		fmt.Println("Error generating access JWT token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update the access token in the database
	err = updateAccessToken(db, dbTimeout, user.ID, accessToken, 30*time.Minute)
	if err != nil {
		fmt.Println("Error updating access token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	response := map[string]interface{}{"message": "Token refreshed successfully", "accessToken": accessToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
