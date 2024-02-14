// main.go
package main

import (
	"backend-project/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

// HelloWorldHandler returns a simple "Hello, World!" message, helps ensures server loads
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello, World!"}`)
}

func main() {
	// Hardcoded database connection details for testing
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	// Construct data source name
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Print dataSourceName for debugging
	log.Println("DataSourceName:", dataSourceName)

	// Initialize database connection
	db, err := sql.Open("mysql", dataSourceName)
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

	// Router initialization
	router := mux.NewRouter()

	// Registering API endpoints

	// Registration endpoint (no authentication required)
	router.HandleFunc("/register", RegisterHandler).Methods("POST")

	// VerifyPin endpoint
	router.HandleFunc("/verify-pin", VerifyPinHandler).Methods("POST")

	// Login endpoint (no authentication required)
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Logout endpoint (requires authentication)
	router.Handle("/logout", validateAccessToken(http.HandlerFunc(LogoutHandler))).Methods("POST")

	// Profile endpoint (requires authentication)
	router.Handle("/profile", validateAccessToken(http.HandlerFunc(ProfileHandler))).Methods("GET")

	// Ticket endpoints

	// Create ticket endpoint
	router.HandleFunc("/tickets", CreateTicketHandler).Methods("POST")

	// Add conversation to ticket endpoint
	router.HandleFunc("/tickets/{ticketID}/conversation", AddConversationHandler).Methods("POST")

	// Get all tickets endpoint
	router.HandleFunc("/tickets", GetTicketsHandler).Methods("GET")

	// Get ticket by ID endpoint
	router.HandleFunc("/tickets/{ticketID}", GetTicketByIDHandler).Methods("GET")

	// Close ticket endpoint
	router.HandleFunc("/tickets/{ticketID}", CloseTicketHandler).Methods("DELETE")

	// Admin endpoints

	// View all tickets (requires admin privilege)
	router.Handle("/admin/tickets", validateAdminAccess(http.HandlerFunc(ViewAllTicketsHandler))).Methods("GET")

	// Get ticket by ID for admin endpoint
	router.Handle("/admin/tickets/{ticketID}", validateAdminAccess(http.HandlerFunc(AdminGetTicketByIDHandler))).Methods("GET")

	// Add conversation to ticket for admin endpoint
	router.Handle("/admin/tickets/{ticketID}/conversation", validateAdminAccess(http.HandlerFunc(AdminAddConversationHandler))).Methods("POST")

	// Token refreshing endpoint
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
