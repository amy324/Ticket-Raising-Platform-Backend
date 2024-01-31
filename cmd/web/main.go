package main

import (
	//"database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend-project/data"

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create the user in the database
	userID, err := user.Create()
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"message": "User registered successfully", "userID": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LoginHandler
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
	// Generate access token with 30 minutes expiry
	accessToken, err := generateJWT(user, os.Getenv("JWT_ACCESS_KEY"), 30*time.Minute)
	if err != nil {
		fmt.Println("Error generating access JWT token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate refresh token with 30 days expiry
	refreshToken, err := generateJWT(user, os.Getenv("JWT_REFRESH_KEY"), 30*24*time.Hour)
	if err != nil {
		fmt.Println("Error generating refresh JWT token:", err)
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

	// Include both tokens in the response
	response := map[string]interface{}{
		"message":      "Login successful",
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Log successful login
	fmt.Printf("User %s successfully logged in\n", user.Email)
}

// Function to generate JWT token with expiration time
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

	// Extract user ID from the request context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)
		return
	}

	// You can now use the userID for further processing
	response := map[string]interface{}{"message": "Protected endpoint", "userID": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ...
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// You can perform any additional cleanup or token invalidation here

	// Return a response without a token to simulate logout
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

	// Logout endpoint (no authentication required)
	router.HandleFunc("/logout", LogoutHandler).Methods("POST") // <-- Fix: Use LogoutHandler

	// Hello, World! endpoint (no authentication required)
	router.HandleFunc("/", HelloWorldHandler).Methods("GET")

	// Start the server
	port := 8080
	log.Printf("Server started on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
