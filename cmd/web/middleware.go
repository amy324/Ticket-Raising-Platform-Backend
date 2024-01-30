package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

// Middleware to validate JWT token
func validateToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract the token from the Authorization header
        tokenString := r.Header.Get("Authorization")

        // Check if the token is present
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Parse the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Check the signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            // Replace "your_secret_key" with your actual secret key
            return []byte("your_secret_key"), nil
        })

        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check if the token is valid
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Extract user ID from the token
            userID, err := strconv.Atoi(claims["sub"].(string))
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Attach user ID to the request context
            ctx := context.WithValue(r.Context(), "userID", userID)
            r = r.WithContext(ctx)

            // Call the next handler
            next.ServeHTTP(w, r)
        } else {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
    })
}
