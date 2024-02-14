// admin.go

package main

import (
	"backend-project/data" 
	"log"
	"net/http"
	"strings"
)

// Middleware to validate admin access
func validateAdminAccess(next http.Handler) http.Handler {
	// Return an anonymous function as an http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the access token from the Authorization header
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "Access token is required", http.StatusBadRequest)
			return
		}

		// Check if the token starts with the "Bearer " prefix
		if strings.HasPrefix(accessToken, "Bearer ") {
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		}

		// Check if the access token is expired
		if isTokenExpired(accessToken) {
			http.Error(w, "Access token has expired", http.StatusUnauthorized)
			return
		}

		// Extract user ID associated with the access token
		userID, err := data.GetUserIDByAccessToken(accessToken)
		if err != nil {
			log.Println("Failed to extract user ID:", err)
			http.Error(w, "Failed to extract user ID", http.StatusInternalServerError)
			return
		}

		// Retrieve user information by user ID
		user, err := data.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
			return
		}

		// Check if the user is an admin
		if user.IsAdmin != 1 {
			http.Error(w, "Access denied. Admin privilege required.", http.StatusForbidden)
			return
		}

		// Call the next handler if admin access is validated
		next.ServeHTTP(w, r)
	})
}
