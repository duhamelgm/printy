package db

import (
	"database/sql"
	"fmt"
	"time"
)

// CreatePrint creates a new print record
func (d *Database) CreatePrint(print *Print) error {
	query := `
		INSERT INTO prints (ticket_id, created_at, updated_at)
		VALUES (?, ?, ?)`

	result, err := d.db.Exec(query, print.TicketID, print.CreatedAt, print.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create print: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	print.ID = int(id)
	return nil
}

// GetPrintByID retrieves a print by ID
func (d *Database) GetPrintByID(id int) (*Print, error) {
	query := `SELECT id, ticket_id, created_at, updated_at FROM prints WHERE id = ?`

	print := &Print{}
	err := d.db.QueryRow(query, id).Scan(
		&print.ID, &print.TicketID, &print.CreatedAt, &print.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("print not found")
		}
		return nil, fmt.Errorf("failed to get print: %v", err)
	}

	return print, nil
}

// GetPrintsByTicketID retrieves all prints for a specific ticket
func (d *Database) GetPrintsByTicketID(ticketID int) ([]Print, error) {
	query := `SELECT id, ticket_id, created_at, updated_at FROM prints WHERE ticket_id = ? ORDER BY created_at DESC`

	rows, err := d.db.Query(query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query prints: %v", err)
	}
	defer rows.Close()

	var prints []Print
	for rows.Next() {
		var print Print
		err := rows.Scan(
			&print.ID, &print.TicketID, &print.CreatedAt, &print.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan print: %v", err)
		}
		prints = append(prints, print)
	}

	return prints, nil
}

// GetAllPrints retrieves all prints
func (d *Database) GetAllPrints() ([]Print, error) {
	query := `SELECT id, ticket_id, created_at, updated_at FROM prints ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query prints: %v", err)
	}
	defer rows.Close()

	var prints []Print
	for rows.Next() {
		var print Print
		err := rows.Scan(
			&print.ID, &print.TicketID, &print.CreatedAt, &print.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan print: %v", err)
		}
		prints = append(prints, print)
	}

	return prints, nil
}

// GetTicketWithPrints retrieves a ticket with all its associated prints
func (d *Database) GetTicketWithPrints(ticketID int) (*TicketWithPrints, error) {
	// Get the ticket
	ticket, err := d.GetTicketByID(ticketID)
	if err != nil {
		return nil, err
	}

	// Get all prints for this ticket
	prints, err := d.GetPrintsByTicketID(ticketID)
	if err != nil {
		return nil, err
	}

	return &TicketWithPrints{
		Ticket: *ticket,
		Prints: prints,
	}, nil
}

// GetPrintsByDateRange retrieves prints within a date range
func (d *Database) GetPrintsByDateRange(startDate, endDate time.Time) ([]Print, error) {
	query := `SELECT id, ticket_id, created_at, updated_at FROM prints WHERE created_at BETWEEN ? AND ? ORDER BY created_at DESC`

	rows, err := d.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query prints by date range: %v", err)
	}
	defer rows.Close()

	var prints []Print
	for rows.Next() {
		var print Print
		err := rows.Scan(
			&print.ID, &print.TicketID, &print.CreatedAt, &print.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan print: %v", err)
		}
		prints = append(prints, print)
	}

	return prints, nil
}

// GetPrintStats returns statistics about prints
func (d *Database) GetPrintStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total prints count
	var totalPrints int
	err := d.db.QueryRow("SELECT COUNT(*) FROM prints").Scan(&totalPrints)
	if err != nil {
		return nil, fmt.Errorf("failed to get total prints count: %v", err)
	}
	stats["total_prints"] = totalPrints

	// Total tickets count
	var totalTickets int
	err = d.db.QueryRow("SELECT COUNT(*) FROM tickets").Scan(&totalTickets)
	if err != nil {
		return nil, fmt.Errorf("failed to get total tickets count: %v", err)
	}
	stats["total_tickets"] = totalTickets

	// Prints today
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	var printsToday int
	err = d.db.QueryRow("SELECT COUNT(*) FROM prints WHERE created_at BETWEEN ? AND ?", today, tomorrow).Scan(&printsToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get prints today count: %v", err)
	}
	stats["prints_today"] = printsToday

	// Most printed ticket
	var mostPrintedTicketID int
	var mostPrintedCount int
	err = d.db.QueryRow(`
		SELECT ticket_id, COUNT(*) as print_count 
		FROM prints 
		GROUP BY ticket_id 
		ORDER BY print_count DESC 
		LIMIT 1
	`).Scan(&mostPrintedTicketID, &mostPrintedCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get most printed ticket: %v", err)
	}

	if err == nil {
		stats["most_printed_ticket_id"] = mostPrintedTicketID
		stats["most_printed_count"] = mostPrintedCount
	}

	return stats, nil
}

// DeletePrint deletes a print record
func (d *Database) DeletePrint(id int) error {
	query := `DELETE FROM prints WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete print: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("print not found")
	}

	return nil
}
