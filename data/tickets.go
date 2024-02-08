package data

import (
	"context"
	"log"

	//"database/sql"
	"time"
)

// Ticket represents the structure of a ticket.
type Ticket struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userId"`
	Email      string    `json:"email"`
	Subject    string    `json:"subject"`
	Issue      string    `json:"issue"`
	Status     string    `json:"status"`
	DateOpened time.Time `json:"dateOpened"`
}

// Conversation represents a message within a ticket conversation.
type Conversation struct {
	ID            int       `json:"id"`
	TicketID      int       `json:"ticketId"`
	Sender        string    `json:"sender"`
	Message       string    `json:"message"`
	MessageSentAt time.Time `json:"messageSentAt"`
}

// CreateTicket creates a new ticket in the database and returns its ID.
func CreateTicket(userID int, subject, issue string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Retrieve the user's email by userID
	userEmail, err := GetUserEmailByID(userID)
	if err != nil {
		log.Printf("Error retrieving user email: %v", err)
		return 0, err
	}

	// Prepare the SQL statement to insert a new ticket
	stmt := `
        INSERT INTO tickets (userId, email, subject, issue, status, dateOpened)
        VALUES (?, ?, ?, ?, ?, NOW())`

	// Execute the SQL statement
	result, err := db.ExecContext(ctx, stmt, userID, userEmail, subject, issue, "open")
	if err != nil {
		log.Printf("Error inserting ticket into database: %v", err)
		return 0, err
	}

	// Retrieve the ID of the newly inserted ticket
	ticketID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last inserted ID: %v", err)
		return 0, err
	}

	// Insert the initial conversation for the ticket
	_, err = addInitialConversation(int(ticketID), "operator", "have you tried turning it on and off again?")
	if err != nil {
		log.Printf("Error adding initial conversation: %v", err)
		return 0, err
	}

	log.Printf("Ticket created successfully with ID: %d", ticketID)

	return int(ticketID), nil
}

// addInitialConversation inserts the initial conversation message for a ticket.
func addInitialConversation(ticketID int, sender, message string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Prepare the SQL statement to insert the initial conversation
	stmt := `
        INSERT INTO conversations (ticketId, sender, message, messageSentAt)
        VALUES (?, ?, ?, NOW())`

	// Execute the SQL statement
	result, err := db.ExecContext(ctx, stmt, ticketID, sender, message)
	if err != nil {
		return 0, err
	}

	// Retrieve the ID of the newly inserted conversation
	conversationID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(conversationID), nil
}

// AddConversation adds a conversation to a ticket in the database.
func AddConversation(ticketID int64, sender, message string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, "INSERT INTO ticket_conversations (ticketId, sender, message, message_sent_at) VALUES (?, ?, ?, ?)",
		ticketID, sender, message, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetTicketsByUserID retrieves all tickets for a given user ID.
func GetTicketsByUserID(userID int64) ([]Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM tickets WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(&ticket.ID, &ticket.UserID, &ticket.Email, &ticket.Subject, &ticket.Status, &ticket.DateOpened)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

// GetTicketByID retrieves a ticket by its ID.
func GetTicketByID(ticketID int64) (Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var ticket Ticket
	err := db.QueryRowContext(ctx, "SELECT * FROM tickets WHERE id = ?", ticketID).
		Scan(&ticket.ID, &ticket.UserID, &ticket.Email, &ticket.Subject, &ticket.Status, &ticket.DateOpened)
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

// CloseTicket closes a ticket by updating its status to "closed".
func CloseTicket(ticketID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, "UPDATE tickets SET status = ? WHERE id = ?", "closed", ticketID)
	return err
}
