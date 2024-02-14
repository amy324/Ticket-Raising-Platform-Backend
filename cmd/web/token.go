// token.go

package main

import (
	"backend-project/data"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// generateTokens generates access and refresh tokens for the given user
func generateTokens(user *data.User) (string, string, error) {
	// Generate access token with 30 minutes expiry
	accessToken, err := generateAuthJWT(user, os.Getenv("JWT_ACCESS_KEY"), 30*time.Minute)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with 30 days expiry
	refreshToken, err := generateAuthJWT(user, os.Getenv("JWT_REFRESH_KEY"), 30*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateAuthJWT generates a JWT token with the given user information, secret key, and expiration time
func generateAuthJWT(user *data.User, secretKey string, expirationTime time.Duration) (string, error) {
	// Set the expiration time for the token
	expiration := time.Now().Add(expirationTime)

	// Create the JWT claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(user.ID),
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the provided secret key
	key := []byte(secretKey)
	if len(key) == 0 {
		log.Fatal("JWT secret key is not set")
	}

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// validateAccessToken validates the access token provided in the request header
func validateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the access token from the Authorization header
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "Access token is required", http.StatusBadRequest)
			return
		}

		// Check if the token starts with the "Bearer " prefix
		if strings.HasPrefix(accessToken, "Bearer ") {
			// Remove the "Bearer " prefix from the token
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		}

		// Log the received access token for debugging
		fmt.Println("Received access token:", accessToken)

		// Validate the access token expiration
		if isTokenExpired(accessToken) {
			http.Error(w, "Access token has expired", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// validateRefreshJWT validates the refresh token and retrieves user information
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

// isTokenExpired checks if the access token is expired
func isTokenExpired(accessToken string) bool {
	// Get user ID associated with the access token
	userID, err := data.GetUserIDByAccessToken(accessToken)
	if err != nil {
		fmt.Println("Error fetching user ID from access token:", err)
		return true // Consider expired if unable to fetch user ID
	}

	// Get expiration time from the database
	expirationTime, err := data.GetAccessTokenExpirationTime(userID)
	if err != nil {
		fmt.Println("Error fetching token expiration time:", err)
		return true // Consider expired if unable to fetch expiration time
	}

	fmt.Println("Expiration time from database:", expirationTime)

	// Check if the current time is after the expiration time
	currentTime := time.Now()
	fmt.Println("Current time:", currentTime)
	if currentTime.After(expirationTime) {
		fmt.Println("Access token has expired")
		return true
	}

	return false
}
