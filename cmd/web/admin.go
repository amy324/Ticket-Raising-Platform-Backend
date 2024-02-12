package main

import (
	"backend-project/data"
	"log"
	"net/http"
	"strings"
)
// Middleware to validate admin access
func validateAdminAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "Access token is required", http.StatusBadRequest)
			return
		}

		if strings.HasPrefix(accessToken, "Bearer ") {
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		}

		if isTokenExpired(accessToken) {
			http.Error(w, "Access token has expired", http.StatusUnauthorized)
			return
		}

		userID, err := data.GetUserIDByAccessToken(accessToken)
		if err != nil {
			log.Println("Failed to extract user ID:", err)
			http.Error(w, "Failed to extract user ID", http.StatusInternalServerError)
			return
		}

		user, err := data.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
			return
		}

		if user.IsAdmin != 1 {
			http.Error(w, "Access denied. Admin privilege required.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
