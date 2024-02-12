package data

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// User structure
type User struct {
	ID         int
	Email      string
	FirstName  string
	LastName   string
	Password   string
	PinNumber  string
	UserActive int
	IsAdmin    int
	RefreshJWT string
}

// AccessToken structure
type AccessToken struct {
	ID        int
	UserID    int
	Email     string
	AccessJWT string
}

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

