package db

import (
	"encoding/json"
	"time"
)

// Ticket represents a ticket in the database
type Ticket struct {
	ID        int       `json:"id" db:"id"`
	RefID     string    `json:"ref_id" db:"ref_id"`
	Title     string    `json:"title" db:"title"`
	Priority  int       `json:"priority" db:"priority"`
	Cooldown  int       `json:"cooldown" db:"cooldown"` // Cooldown in seconds
	Weekdays  string    `json:"weekdays" db:"weekdays"` // Weekdays as JSON array string (e.g., ["WeekEnd", "WeekDay"])
	Assignee  string    `json:"assignee" db:"assignee"` // Assignee name from Notion user
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Print represents a print job in the database
type Print struct {
	ID        int       `json:"id" db:"id"`
	TicketID  int       `json:"ticket_id" db:"ticket_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TicketWithPrints represents a ticket with its associated prints
type TicketWithPrints struct {
	Ticket Ticket  `json:"ticket"`
	Prints []Print `json:"prints"`
}

// GetWeekdaysAsArray returns the weekdays as a slice of strings
func (t *Ticket) GetWeekdaysAsArray() ([]string, error) {
	if t.Weekdays == "" {
		return []string{}, nil
	}

	var weekdays []string
	err := json.Unmarshal([]byte(t.Weekdays), &weekdays)
	if err != nil {
		// If JSON parsing fails, try comma-separated parsing for backward compatibility
		if t.Weekdays != "" {
			weekdays = []string{t.Weekdays}
		}
	}
	return weekdays, nil
}

// SetWeekdaysFromArray sets the weekdays from a slice of strings
func (t *Ticket) SetWeekdaysFromArray(weekdays []string) error {
	if len(weekdays) == 0 {
		t.Weekdays = "[]"
		return nil
	}

	jsonData, err := json.Marshal(weekdays)
	if err != nil {
		return err
	}

	t.Weekdays = string(jsonData)
	return nil
}

// HasWeekday checks if the ticket has a specific weekday type
func (t *Ticket) HasWeekday(weekday string) (bool, error) {
	weekdays, err := t.GetWeekdaysAsArray()
	if err != nil {
		return false, err
	}

	for _, w := range weekdays {
		if w == weekday {
			return true, nil
		}
	}
	return false, nil
}
