package data

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// User structure represents the attributes of a user in the system.
type User struct {
	ID         int    // Unique identifier for the user
	Email      string // Email address of the user
	FirstName  string // First name of the user
	LastName   string // Last name of the user
	Password   string // Hashed password of the user
	PinNumber  string // PIN number associated with the user
	UserActive int    // Flag indicating whether the user account is active (1) or not (0)
	IsAdmin    int    // Flag indicating whether the user is an administrator (1) or not (0)
	RefreshJWT string // Refresh JSON Web Token (JWT) for the user
}

// AccessToken structure represents the access token associated with a user.
type AccessToken struct {
	ID        int    // Unique identifier for the access token
	UserID    int    // ID of the user associated with the access token
	Email     string // Email address of the user
	AccessJWT string // Access JSON Web Token (JWT)
}

// Ticket represents the structure of a ticket in the system.
type Ticket struct {
	ID         int64     `json:"id"`          // Unique identifier for the ticket
	UserID     int64     `json:"userId"`      // ID of the user who opened the ticket
	Email      string    `json:"email"`       // Email address of the user who opened the ticket
	Subject    string    `json:"subject"`     // Subject of the ticket
	Issue      string    `json:"issue"`       // Description of the issue
	Status     string    `json:"status"`      // Status of the ticket (e.g., open, closed)
	DateOpened time.Time `json:"dateOpened"`  // Date and time when the ticket was opened
}

// Conversation represents a message within a ticket conversation.
type Conversation struct {
	ID            int64     `json:"id"`            // Unique identifier for the conversation message
	TicketID      int64     `json:"ticketId"`      // ID of the ticket associated with the conversation
	Sender        string    `json:"sender"`        // Sender of the message
	Message       string    `json:"message"`       // Content of the message
	MessageSentAt time.Time `json:"messageSentAt"` // Date and time when the message was sent
}
