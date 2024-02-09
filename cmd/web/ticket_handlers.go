package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"backend-project/data"

	"github.com/gorilla/mux"
)

func CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Creating ticket...")

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

	// Extract the user ID associated with the access token
	userID, err := data.GetUserIDByAccessToken(accessToken)
	if err != nil {
		log.Println("Failed to extract user ID:", err)
		http.Error(w, "Failed to extract user ID", http.StatusInternalServerError)
		return
	}

	// Parse request body
	var ticketData struct {
		Subject string `json:"subject"`
		Issue   string `json:"issue"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ticketData); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Create ticket
	ticketID, err := data.CreateTicket(userID, ticketData.Subject, ticketData.Issue)

	if err != nil {
		log.Println("Error creating ticket:", err)
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	// Respond with ticket ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{ TicketID int }{TicketID: ticketID})
}

// AddConversationHandler handles requests to add a conversation to a ticket.
func AddConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Adding conversation..")

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
	// Extract ticketID from request URL
	params := mux.Vars(r)
	ticketID, err := strconv.ParseInt(params["ticketID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var conversation struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&conversation); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Add conversation to ticket
	_, err = data.AddConversation(ticketID, "user", conversation.Message)
	if err != nil {
		http.Error(w, "Failed to add conversation to ticket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
// GetTicketsHandler retrieves tickets for the user associated with the access token.
func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
    // Log the start of the handler
    log.Println("Getting all tickets...")

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

    // Get user ID as int64
    userID, err := data.GetUserIDByAccessTokenInt64(accessToken)
    if err != nil {
        log.Printf("Failed to retrieve user ID: %v", err)
        http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
        return
    }

    // Get tickets for user
    tickets, err := data.GetTicketsByUserID(userID)
    if err != nil {
        log.Printf("Failed to retrieve tickets: %v", err)
        http.Error(w, "Failed to retrieve tickets", http.StatusInternalServerError)
        return
    }

    // Respond with tickets
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tickets)
}




// GetTicketByIDHandler handles requests to retrieve a specific ticket by its ID.
// GetTicketByIDHandler handles requests to retrieve a specific ticket by its ID along with its conversations.
func GetTicketByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Log the start of the handler
    log.Println("Getting ticket by ID...")

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

    // Extract ticketID from request URL
    params := mux.Vars(r)
    ticketID, err := strconv.ParseInt(params["ticketID"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
        return
    }

    // Get ticket details by ID
    ticket, err := data.GetTicketByID(ticketID)
    if err != nil {
        http.Error(w, "Failed to retrieve ticket", http.StatusInternalServerError)
        return
    }

    // Get conversations for the ticket
    conversations, err := data.GetConversationsByTicketID(ticketID)
    if err != nil {
        http.Error(w, "Failed to retrieve conversations", http.StatusInternalServerError)
        return
    }

    // Combine ticket and conversations into a struct
    type TicketWithConversations struct {
        Ticket        data.Ticket        `json:"ticket"`
        Conversations []data.Conversation `json:"conversations"`
    }

    // Create the combined data
    ticketWithConversations := TicketWithConversations{
        Ticket:        ticket,
        Conversations: conversations,
    }

    // Respond with combined data
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ticketWithConversations)
}

// CloseTicketHandler handles requests to close a ticket.
func CloseTicketHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Closing ticket...")

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
	// Extract ticketID from request URL
	params := mux.Vars(r)
	ticketID, err := strconv.ParseInt(params["ticketID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	// Close ticket
	err = data.CloseTicket(ticketID)
	if err != nil {
		http.Error(w, "Failed to close ticket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
