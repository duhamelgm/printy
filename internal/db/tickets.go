package db

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateTicket creates a new ticket
func (d *Database) CreateTicket(ticket *Ticket) error {
	query := `
		INSERT INTO tickets (ref_id, title, priority, cooldown, weekdays, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, ticket.RefID, ticket.Title, ticket.Priority, ticket.Cooldown, ticket.Weekdays, ticket.CreatedAt, ticket.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	ticket.ID = int(id)
	return nil
}

// GetTicketByID retrieves a ticket by ID
func (d *Database) GetTicketByID(id int) (*Ticket, error) {
	query := `SELECT id, ref_id, title, priority, cooldown, weekdays, created_at, updated_at FROM tickets WHERE id = ?`

	ticket := &Ticket{}
	err := d.db.QueryRow(query, id).Scan(
		&ticket.ID, &ticket.RefID, &ticket.Title, &ticket.Priority,
		&ticket.Cooldown, &ticket.Weekdays, &ticket.CreatedAt, &ticket.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, fmt.Errorf("failed to get ticket: %v", err)
	}

	return ticket, nil
}

// GetTicketByRefID retrieves a ticket by reference ID
func (d *Database) GetTicketByRefID(refID string) (*Ticket, error) {
	query := `SELECT id, ref_id, title, priority, cooldown, weekdays, created_at, updated_at FROM tickets WHERE ref_id = ?`

	ticket := &Ticket{}
	err := d.db.QueryRow(query, refID).Scan(
		&ticket.ID, &ticket.RefID, &ticket.Title, &ticket.Priority,
		&ticket.Cooldown, &ticket.Weekdays, &ticket.CreatedAt, &ticket.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, fmt.Errorf("failed to get ticket: %v", err)
	}

	return ticket, nil
}

// GetAllTickets retrieves all tickets
func (d *Database) GetAllTickets() ([]Ticket, error) {
	query := `SELECT id, ref_id, title, priority, cooldown, weekdays, created_at, updated_at FROM tickets ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %v", err)
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(
			&ticket.ID, &ticket.RefID, &ticket.Title, &ticket.Priority,
			&ticket.Cooldown, &ticket.Weekdays, &ticket.CreatedAt, &ticket.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %v", err)
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// UpdateTicket updates an existing ticket
func (d *Database) UpdateTicket(ticket *Ticket) error {
	query := `
		UPDATE tickets 
		SET ref_id = ?, title = ?, priority = ?, cooldown = ?, weekdays = ?, updated_at = ?
		WHERE id = ?`

	result, err := d.db.Exec(query, ticket.RefID, ticket.Title, ticket.Priority, ticket.Cooldown, ticket.Weekdays, ticket.UpdatedAt, ticket.ID)
	if err != nil {
		return fmt.Errorf("failed to update ticket: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}

// DeleteTicket deletes a ticket and all associated prints
func (d *Database) DeleteTicket(id int) error {
	query := `DELETE FROM tickets WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}

// GetTicketsByPriority retrieves tickets by priority level
func (d *Database) GetTicketsByPriority(priority int) ([]Ticket, error) {
	query := `SELECT id, ref_id, title, priority, cooldown, weekdays, created_at, updated_at FROM tickets WHERE priority = ? ORDER BY created_at DESC`

	rows, err := d.db.Query(query, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets by priority: %v", err)
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(
			&ticket.ID, &ticket.RefID, &ticket.Title, &ticket.Priority,
			&ticket.Cooldown, &ticket.Weekdays, &ticket.CreatedAt, &ticket.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %v", err)
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// IsTicketInCooldown checks if a ticket is still in cooldown period
func (d *Database) IsTicketInCooldown(ticketID int) (bool, error) {
	ticket, err := d.GetTicketByID(ticketID)
	if err != nil {
		return false, err
	}

	if ticket.Cooldown <= 0 {
		return false, nil
	}

	// Get the last print for this ticket
	query := `SELECT created_at FROM prints WHERE ticket_id = ? ORDER BY created_at DESC LIMIT 1`
	var lastPrintTime time.Time
	err = d.db.QueryRow(query, ticketID).Scan(&lastPrintTime)

	if err == sql.ErrNoRows {
		// No prints yet, not in cooldown
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get last print time: %v", err)
	}

	// Check if cooldown period has passed
	cooldownDuration := time.Duration(ticket.Cooldown) * time.Second
	return time.Since(lastPrintTime) < cooldownDuration, nil
}
