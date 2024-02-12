package data

import "context"

// GetTickets retrieves all tickets from the database.
func GetTickets() ([]Ticket, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    rows, err := db.QueryContext(ctx, "SELECT _id, userId, email, subject, issue, status, dateOpened FROM tickets")
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
