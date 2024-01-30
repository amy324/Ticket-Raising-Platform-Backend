package main

import (
	//"database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend-project/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Temporarily storing the user's password without hashing for testing
	// Remove or comment out the following lines to revert to hashed passwords
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	// if err != nil {
	// 	log.Println("Error hashing the password:", err)
	// 	http.Error(w, "Error hashing the password", http.StatusInternalServerError)
	// 	return
	// }
	// user.Password = string(hashedPassword)

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
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println("Received login request for user:", credentials.Email)

	// Authenticate the user
	user, err := data.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		log.Println("Error authenticating user:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Validate the retrieved user's password
	matches, err := user.PasswordMatches(credentials.Password)
	if err != nil {
		log.Println("Error checking password:", err)
		http.Error(w, "Error checking password", http.StatusInternalServerError)
		return
	}

	log.Println("Password matches:", matches)

	if !matches {
		log.Println("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// You can generate a JWT or create a session here for further authentication

	response := map[string]interface{}{"message": "Login successful", "user": user}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello, World!"}`)
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

	router := mux.NewRouter()

	// Registration endpoint
	router.HandleFunc("/register", RegisterHandler).Methods("POST")

	// Login endpoint
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Hello, World! endpoint
	router.HandleFunc("/", HelloWorldHandler).Methods("GET")

	// Serve static files (like favicon.ico)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start the server
	port := 8080
	log.Printf("Server started on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
