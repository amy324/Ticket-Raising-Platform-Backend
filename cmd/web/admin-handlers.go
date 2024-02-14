// admin_handlers.go

package main

import (
	"backend-project/data" 
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" 
)

// ViewAllTicketsHandler displays all existing tickets without fetching their associated messages
func ViewAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Viewing all tickets...")

	// Fetch all tickets from the database
	tickets, err := data.GetTickets()
	if err != nil {
		http.Error(w, "Failed to fetch tickets.", http.StatusInternalServerError)
		return
	}

	// Serialize tickets to JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

// AdminGetTicketByIDHandler handles requests to retrieve a specific ticket by its ID along with its conversations.
func AdminGetTicketByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Getting ticket by ID...")

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
		Ticket        data.Ticket         `json:"ticket"`
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

// AdminAddConversationHandler adds a conversation to a ticket for admin users
func AdminAddConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Log the start of the handler
	log.Println("Adding conversation...")

	// Extract the ticket ID from the request URL parameters
	params := mux.Vars(r)
	ticketID, err := strconv.ParseInt(params["ticketID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the conversation message
	var conversation struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&conversation); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Ensure that the sender is always "operator" for admin users
	sender := "operator"

	// Add the conversation to the database
	conversationID, err := data.AddConversation(ticketID, sender, conversation.Message)
	if err != nil {
		log.Println("Failed to add conversation to ticket:", err)
		http.Error(w, "Failed to add conversation to ticket", http.StatusInternalServerError)
		return
	}

	// Respond with the conversation ID
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"conversationID": conversationID,
	}
	json.NewEncoder(w).Encode(response)
}
