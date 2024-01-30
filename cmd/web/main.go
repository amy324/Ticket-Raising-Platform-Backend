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

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		fmt.Println("Error generating JWT token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Include the token in the response
	response := map[string]interface{}{"message": "Login successful", "user": user, "token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Log successful login
	fmt.Printf("User %s successfully logged in\n", user.Email)
}

// Function to generate JWT token
func generateJWT(user *data.User) (string, error) {
	// Set the expiration time for the token (you can customize this)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(user.ID),
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key (replace with your own secret key)
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}
	signedToken, err := token.SignedString(secretKey)
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

	// Hello, World! endpoint (no authentication required)
	router.HandleFunc("/", HelloWorldHandler).Methods("GET")

	// Start the server
	port := 8080
	log.Printf("Server started on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
