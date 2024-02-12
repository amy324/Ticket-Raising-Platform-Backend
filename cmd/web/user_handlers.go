package main

import (
	"backend-project/data"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	exists, err := data.UserExists(user.Email)
	if err != nil {
		log.Println("Error checking user existence:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Generate a pin number for verification
	pinNumber, err := data.GeneratePinNumber()
	if err != nil {
		log.Println("Error generating pin number:", err)
		http.Error(w, "Error generating pin number", http.StatusInternalServerError)
		return
	}

	// Set the generated pin number for the user
	user.PinNumber = pinNumber

	// Create the user in the database
	userID, err := user.Create()
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Retrieve the PIN from the database to ensure consistency
	savedPin, err := data.GetPinByEmail(user.Email)
	if err != nil {
		log.Println("Error retrieving PIN from the database:", err)
		http.Error(w, "Error retrieving PIN from the database", http.StatusInternalServerError)
		return
	}

	// Send the PIN via email
	subject := "Verification Code"
	body := fmt.Sprintf("Verification code for user %s: %s", user.Email, savedPin)
	err = sendPinByEmail(user.Email, subject, body)
	if err != nil {
		log.Println("Error sending PIN via email:", err)
		http.Error(w, "Error sending PIN via email", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"message": "User registered successfully", "userID": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// VerifyPinHandler handles PIN verification
func VerifyPinHandler(w http.ResponseWriter, r *http.Request) {
	var pinVerification struct {
		Email string `json:"email"`
		Pin   string `json:"pin"`
	}

	err := json.NewDecoder(r.Body).Decode(&pinVerification)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the user by email
	user, err := data.GetUserByEmail(pinVerification.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Retrieve the PIN for the user from the database
	savedPin, err := data.GetPinByEmail(pinVerification.Email)
	if err != nil {
		http.Error(w, "Error retrieving PIN from the database", http.StatusInternalServerError)
		return
	}

	// Compare the provided PIN with the one retrieved from the database
	if savedPin != pinVerification.Pin {
		http.Error(w, "Invalid PIN", http.StatusUnauthorized)
		return
	}

	// Update the pin_number field to indicate verification
	user.PinNumber = "N/A - verified" // Set the new value for pin_number
	if err := user.UpdatePinAfterVerification(); err != nil {
		http.Error(w, "Error updating PIN after verification", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{"message": "PIN verified successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println("Received login request for user:", credentials.Email)

	// Authenticate the user
	user, err := data.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		if errors.Is(err, data.ErrUserNotFound) {
			// If the user is not found, return a specific response
			http.Error(w, "User does not exist or has not been activated. Please try re-registering your account", http.StatusNotFound)
			return
		} else {
			// For other authentication errors, return generic unauthorized response
			fmt.Println("Error authenticating user:", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	// Check if the user is active
	if user.UserActive != 1 {
		// If the user is not active, return a specific response
		http.Error(w, "User does not exist or has not been activated. Please try re-registering your account", http.StatusForbidden)
		return
	}

	// Insert the access token and update/insert the refresh token into the database
	accessToken, refreshToken, err := generateTokens(user)
	if err != nil {
		fmt.Println("Error generating tokens:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert the access token into the database
	_, err = data.CreateAccessToken(user.ID, user.Email, accessToken)
	if err != nil {
		fmt.Println("Error inserting access token into the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update or insert the refresh token into the database
	err = data.UpdateRefreshToken(user.ID, refreshToken)
	if err != nil {
		fmt.Println("Error updating refresh token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":      "Login successful",
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("User %s successfully logged in\n", user.Email)
}


// LogoutHandler for user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Print the request URL and headers for debugging
	fmt.Println("Request URL:", r.URL)
	fmt.Println("Request Headers:", r.Header)

	// Extract the access token from the Authorization header
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		fmt.Println("No access token provided")
		http.Error(w, "Access token is required", http.StatusBadRequest)
		return
	}

	// Remove the "Bearer " prefix from the token
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	// Retrieve user ID from the database using the access token
	userID, err := data.GetUserIDByAccessToken(accessToken)
	if err != nil {
		fmt.Println("Error retrieving user ID from access tokens table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve the user from the database based on the user ID
	user, err := data.GetUserByID(userID)
	if err != nil {
		fmt.Println("Error retrieving user from the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("User with ID %d is logging out\n", userID)

	// Logout the user using the Logout method defined on the User struct
	if err := user.Logout(); err != nil {
		fmt.Println("Error logging out:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{"message": "Logout successful"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Fetching user profile...")

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

	// Check if the access token is expired
	if isTokenExpired(accessToken) {
		http.Error(w, "Access token has expired", http.StatusUnauthorized)
		return
	}

	// Extract the email associated with the access token
	email, err := data.GetUserEmailByAccessToken(accessToken)
	if err != nil {
		log.Println("Error retrieving user email from access token:", err)
		http.Error(w, "Error retrieving user email", http.StatusInternalServerError)
		return
	}

	// Retrieve the user profile from the database using the email
	user, err := data.GetUserByEmail(email)
	if err != nil {
		log.Println("Error retrieving user profile:", err)
		http.Error(w, "Error retrieving user profile", http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Println("User profile not found")
		http.Error(w, "User profile not found", http.StatusNotFound)
		return
	}

	// Log the retrieved user profile
	log.Println("Retrieved user profile:", user)

	// Respond with the user profile
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
