package data

import (
	"context"
	"database/sql"
	"log"

	//"database/sql"
	"time"
)

// Ticket represents the structure of a ticket.
type Ticket struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	Email      string    `json:"email"`
	Subject    string    `json:"subject"`
	Issue      string    `json:"issue"`
	Status     string    `json:"status"`
	DateOpened time.Time `json:"dateOpened"`
}

// Conversation represents a message within a ticket conversation.
type Conversation struct {
	ID            int64     `json:"id"`
	TicketID      int64     `json:"ticketId"`
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

	rows, err := db.QueryContext(ctx, "SELECT _id, userId, email, subject, issue, status, dateOpened FROM tickets WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(&ticket.ID, &ticket.UserID, &ticket.Email, &ticket.Subject, &ticket.Issue, &ticket.Status, &ticket.DateOpened)
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

// GetTicketsByUserID retrieves all tickets for a given user ID.
func GetTicketByID(userID int64) ([]Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM tickets WHERE userId = ?", userID)
	if err != nil {
		// Log the error
		log.Printf("Error querying tickets by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(&ticket.ID, &ticket.UserID, &ticket.Email, &ticket.Subject, &ticket.Issue, &ticket.Status, &ticket.DateOpened)
		if err != nil {
			// Log the error
			log.Printf("Error scanning ticket row: %v", err)
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	if err := rows.Err(); err != nil {
		// Log the error
		log.Printf("Error iterating over ticket rows: %v", err)
		return nil, err
	}

	return tickets, nil
}

// CloseTicket closes a ticket by updating its status to "closed".
func CloseTicket(ticketID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, "UPDATE tickets SET status = ? WHERE id = ?", "closed", ticketID)
	return err
}

// GetUserIDByAccessTokenInt64 retrieves the user ID associated with the given access token as int64.
func GetUserIDByAccessTokenInt64(accessToken string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var userID int64
	query := `SELECT user_id FROM access_tokens WHERE accessJWT = ?`

	// Log the query being executed
	log.Printf("Executing query to retrieve user ID for access token: %s", accessToken)

	// Execute the query and scan the result
	err := db.QueryRowContext(ctx, query, accessToken).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return 0 if no rows are found
			return 0, nil
		}
		return 0, err
	}

	// Log the retrieved user ID
	log.Printf("Retrieved user ID from database: %d", userID)

	return userID, nil
}
