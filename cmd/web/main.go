// main.go

package main

import (
	"backend-project/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

//helps to ensure server loads properly
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello, World!"}`)
}


func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/backend_db?parseTime=true")
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

	// Define the routes for ticket operations
	router.HandleFunc("/tickets", CreateTicketHandler).Methods("POST")
	router.HandleFunc("/tickets/{ticketID}/conversation", AddConversationHandler).Methods("POST")
	router.HandleFunc("/tickets", GetTicketsHandler).Methods("GET")
	router.HandleFunc("/tickets/{ticketID}", GetTicketByIDHandler).Methods("GET")
	router.HandleFunc("/tickets/{ticketID}", CloseTicketHandler).Methods("DELETE")

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
