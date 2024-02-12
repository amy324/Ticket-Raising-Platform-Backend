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

// helps to ensure server loads properly
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello, World!"}`)
}

func main() {
	// Fetch database connection details from environment variables
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	// Construct data source name
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

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

	// Update your router initialization in the main function
	router := mux.NewRouter()

	// Registration endpoint (no authentication required)
	router.HandleFunc("/register", RegisterHandler).Methods("POST")

	// VerifyP pin endpoint to the router
	router.HandleFunc("/verify-pin", VerifyPinHandler).Methods("POST")

	// Login endpoint (no authentication required)
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Logout endpoint (requires authentication)
	router.Handle("/logout", validateAccessToken(http.HandlerFunc(LogoutHandler))).Methods("POST")

	//Profile endpoint gets users profile
	router.Handle("/profile", validateAccessToken(http.HandlerFunc(ProfileHandler))).Methods("GET")

	// Routes for ticket operations
	router.HandleFunc("/tickets", CreateTicketHandler).Methods("POST")
	router.HandleFunc("/tickets/{ticketID}/conversation", AddConversationHandler).Methods("POST")
	router.HandleFunc("/tickets", GetTicketsHandler).Methods("GET")
	router.HandleFunc("/tickets/{ticketID}", GetTicketByIDHandler).Methods("GET")
	router.HandleFunc("/tickets/{ticketID}", CloseTicketHandler).Methods("DELETE")

	//Routes for admin use
	// View all tickets (requires admin privilege)
	router.Handle("/admin/tickets", validateAdminAccess(http.HandlerFunc(ViewAllTicketsHandler))).Methods("GET")
	router.Handle("/admin/tickets/{ticketID}", validateAdminAccess(http.HandlerFunc(AdminGetTicketByIDHandler))).Methods("GET")
	router.Handle("/admin/tickets/{ticketID}/conversation", validateAdminAccess(http.HandlerFunc(AdminAddConversationHandler))).Methods("POST")

	// Endpoint for token refreshing
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
