package data

import (
	"context"
	//"database/sql"
	"time"
)
// Ticket represents the structure of a ticket.
type Ticket struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"userId"`
	Email         string        `json:"email"`
	Subject       string        `json:"subject"`
	Status        string        `json:"status"`
	Conversations []Conversation `json:"conversations"`
	DateOpened    time.Time     `json:"dateOpened"`
}

// Conversation represents a message within a ticket conversation.
type Conversation struct {
	ID            int64     `json:"id"`
	TicketID      int64     `json:"ticketId"`
	Sender        string    `json:"sender"`
	Message       string    `json:"message"`
	MessageSentAt time.Time `json:"messageSentAt"`
}
// CreateTicket creates a new ticket in the database.
func CreateTicket(userID int64, email, subject string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, "INSERT INTO tickets (user_id, email, subject, status, date_opened) VALUES (?, ?, ?, ?, ?)",
		userID, email, subject, "open", time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// AddConversation adds a conversation to a ticket in the database.
func AddConversation(ticketID int64, sender, message string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := db.ExecContext(ctx, "INSERT INTO ticket_conversations (ticket_id, sender, message, message_sent_at) VALUES (?, ?, ?, ?)",
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