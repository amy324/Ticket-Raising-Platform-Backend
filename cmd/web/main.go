// main.go

package main

import (
	"backend-project/data"
	//"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func validateToken(next http.Handler) http.Handler {
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

		// Retrieve user ID from the access_tokens table using the access token
		userID, err := data.GetUserIDByAccessToken(accessToken)
		if err != nil {
			http.Error(w, "Error retrieving user ID from access tokens table", http.StatusInternalServerError)
			return
		}

		// Log the retrieved user ID for debugging
		fmt.Println("Retrieved user ID from access tokens table:", userID)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a pin number for verification
	pinNumber, err := data.GeneratePinNumber() // Accessing GeneratePinNumber from the data package
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

	// Send the verification PIN to the user's email address
	emailSubject := "Verification PIN"
	emailBody := fmt.Sprintf("Your verification PIN is: %s", pinNumber)
	recipientEmail := user.Email

	// Send email
	if err := sendEmail(recipientEmail, emailSubject, emailBody); err != nil {
		log.Println("Error sending verification email:", err)
		http.Error(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	// Log or print the pin number for verification (replace with actual email sending code)
	log.Println("Verification code for user", user.Email+":", pinNumber)

	response := map[string]interface{}{"message": "User registered successfully", "userID": userID}
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
		fmt.Println("Error authenticating user:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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

func generateTokens(user *data.User) (string, string, error) {
	// Generate access token with 30 minutes expiry
	accessToken, err := generateJWT(user, os.Getenv("JWT_ACCESS_KEY"), 30*time.Minute)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with 30 days expiry
	refreshToken, err := generateJWT(user, os.Getenv("JWT_REFRESH_KEY"), 30*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateJWT(user *data.User, secretKey string, expirationTime time.Duration) (string, error) {
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

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello, World!"}`)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing ProtectedHandler")

	// Extract user information from the request context
	user, ok := r.Context().Value("user").(*data.User)
	if !ok {
		http.Error(w, "Unable to retrieve user information", http.StatusInternalServerError)
		return
	}

	// You can now use the user information for further processing
	response := map[string]interface{}{"message": "Protected endpoint", "user": user}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func refreshAccessToken(w http.ResponseWriter, r *http.Request, refreshToken string, db *sql.DB) {
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

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/backend_db")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	} else {
		log.Println("Connected to the database.")
	}

	// Set the database connection in the data package
	data.SetDB(db)

	// Update your router initialization in the main function
	router := mux.NewRouter()

	// Registration endpoint (no authentication required)
	router.HandleFunc("/register", RegisterHandler).Methods("POST")

	// Login endpoint (no authentication required)
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Protected endpoint (requires authentication)
	router.Handle("/protected", validateToken(http.HandlerFunc(ProtectedHandler))).Methods("GET")

	// Logout endpoint (requires authentication)
	router.Handle("/logout", validateToken(http.HandlerFunc(LogoutHandler))).Methods("POST")

	// Debugging message to ensure /logout endpoint registration
	fmt.Println("Logout endpoint registered successfully.")

	// Add a new endpoint for token refreshing
	router.HandleFunc("/tokens/refresh", func(w http.ResponseWriter, r *http.Request) {
		refreshAccessToken(w, r, r.Header.Get("Authorization"), db)
	}).Methods("POST")

	// Hello, World! endpoint (no authentication required)
	router.HandleFunc("/", HelloWorldHandler).Methods("GET")

	// Start the server
	port := 8080
	log.Printf("Server started on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
