//admin-data.go
package data

import "context"

// GetTickets retrieves all tickets from the database.
func GetTickets() ([]Ticket, error) {
    // Context with timeout to manage database operations
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    // Query to retrieve all tickets
    rows, err := db.QueryContext(ctx, "SELECT id, userId, email, subject, issue, status, dateOpened FROM tickets")
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
