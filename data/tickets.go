//tickets.go
package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// CreateTicket creates a new ticket in the database and returns its ID.
func CreateTicket(userID int, subject, issue string) (int, error) {
    // Context with timeout to manage database operations
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
	_, err = addInitialConversation(int(ticketID), "operator", "We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details")
	if err != nil {
		log.Printf("Error adding initial conversation: %v", err)
		return 0, err
	}

	log.Printf("Ticket created successfully with ID: %d", ticketID)

	return int(ticketID), nil
}

// addInitialConversation inserts the initial conversation message for a ticket.
func addInitialConversation(ticketID int, sender, message string) (int, error) {
    // Context with timeout to manage database operations
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
    // Context with timeout to manage database operations
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Execute the SQL statement to add a conversation
	result, err := db.ExecContext(ctx, "INSERT INTO conversations (ticketId, sender, message, messageSentAt) VALUES (?, ?, ?, ?)",
		ticketID, sender, message, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetTicketsByUserID retrieves all tickets for a given user ID.
func GetTicketsByUserID(userID int64) ([]Ticket, error) {
    // Context with timeout to manage database operations
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query to retrieve tickets by user ID
	rows, err := db.QueryContext(ctx, "SELECT id, userId, email, subject, issue, status, dateOpened FROM tickets WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and populate tickets slice
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

// GetTicketByID retrieves a ticket by its ID from the database.
func GetTicketByID(ticketID int64) (Ticket, error) {
    // Context with timeout to manage database operations
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query to retrieve a ticket by its ID
	var ticket Ticket
	err := db.QueryRowContext(ctx, "SELECT * FROM tickets WHERE id = ?", ticketID).
		Scan(&ticket.ID, &ticket.UserID, &ticket.Email, &ticket.Subject, &ticket.Issue, &ticket.Status, &ticket.DateOpened)
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

// GetConversationsByTicketID retrieves all conversations associated with a ticket ID from the database.
func GetConversationsByTicketID(ticketID int64) ([]Conversation, error) {
    // Context with timeout to manage database operations
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query to retrieve conversations by ticket ID
	rows, err := db.QueryContext(ctx, "SELECT * FROM conversations WHERE ticketId = ?", ticketID)
	if err != nil {
		log.Printf("Error retrieving conversations by ticket ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and populate conversations slice
	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.ID, &conv.TicketID, &conv.Sender, &conv.Message, &conv.MessageSentAt)
		if err != nil {
			log.Printf("Error scanning conversation row: %v", err)
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over conversation rows: %v", err)
		return nil, err
	}

	return conversations, nil
}

// CloseTicket closes a ticket by updating its status to "closed".
func CloseTicket(ticketID int64) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete conversations associated with the ticket
	_, err = tx.Exec("DELETE FROM conversations WHERE ticketId = ?", ticketID)
	if err != nil {
		return err
	}

	// Delete the ticket
	_, err = tx.Exec("DELETE FROM tickets WHERE id = ?", ticketID)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetUserIDByAccessTokenInt64 retrieves the user ID associated with the given access token as int64.
func GetUserIDByAccessTokenInt64(accessToken string) (int64, error) {
    // Context with timeout to manage database operations
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query to retrieve user ID by access token
	var userID int64
	query := `SELECT user_id FROM access_tokens WHERE accessJWT = ?`

	// Execute the query and scan the result
	err := db.QueryRowContext(ctx, query, accessToken).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return 0 if no rows are found
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}

// GetUserIDByTicketID retrieves the user ID associated with a ticket from the database.
func GetUserIDByTicketID(ticketID int64) (int64, error) {
    // Query to retrieve user ID by ticket ID
	var userID int64
	query := "SELECT userId FROM tickets WHERE id = ?"

	err := db.QueryRow(query, ticketID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
