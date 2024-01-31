// handlers.go

package main

import (
	"backend-project/data"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var db *sql.DB
var dbTimeout = 5 * time.Second

// Handler for refreshing the access token
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
	accessToken, err := generateJWT(user, os.Getenv("JWT_ACCESS_KEY"), 30*time.Minute)
	if err != nil {
		fmt.Println("Error generating access JWT token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update the access token in the database
	err = updateAccessToken(db, dbTimeout, user.ID, accessToken)
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

// Function to validate the refresh token and retrieve user information
func validateRefreshJWT(tokenString, secretKey string) (*data.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	// Retrieve user information from the database
	user, err := data.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Function to update the access token in the database
func updateAccessToken(db *sql.DB, timeout time.Duration, userID int, newAccessToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stmt := `
        UPDATE access_tokens
        SET accessJWT = ?
        WHERE user_id = ?`

	_, err := db.ExecContext(ctx, stmt, newAccessToken, userID)
	return err
}
