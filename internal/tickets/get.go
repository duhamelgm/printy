package tickets

import (
	"fmt"
	"printy/internal/db"
	"time"
)

// GetRelevantTickets returns tickets that are relevant for today
func GetRelevantTickets(database *db.Database, count int) ([]db.Ticket, error) {
	// Get all tickets
	allTickets, err := database.GetAllTickets()
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets: %v", err)
	}

	var relevantTickets []db.Ticket
	today := time.Now().Weekday()

	for _, ticket := range allTickets {
		// Check if ticket is relevant for today
		if !isTicketRelevantForToday(ticket, today) {
			continue
		}

		// Check if ticket is in cooldown
		inCooldown, err := database.IsTicketInCooldown(ticket.ID)
		if err != nil {
			continue // Skip if error checking cooldown
		}
		if inCooldown {
			continue // Skip if in cooldown
		}

		relevantTickets = append(relevantTickets, ticket)
	}

	// Sort by priority (higher priority first)
	for i := 0; i < len(relevantTickets)-1; i++ {
		for j := i + 1; j < len(relevantTickets); j++ {
			if relevantTickets[i].Priority < relevantTickets[j].Priority {
				relevantTickets[i], relevantTickets[j] = relevantTickets[j], relevantTickets[i]
			}
		}
	}

	// Limit to requested count
	if len(relevantTickets) > count {
		relevantTickets = relevantTickets[:count]
	}

	return relevantTickets, nil
}

// isTicketRelevantForToday checks if a ticket is relevant for today's weekday
func isTicketRelevantForToday(ticket db.Ticket, today time.Weekday) bool {
	weekdays, err := ticket.GetWeekdaysAsArray()
	if err != nil {
		return false
	}

	// If no weekdays specified, consider it relevant
	if len(weekdays) == 0 {
		return true
	}

	// Check if today matches any of the ticket's weekdays
	for _, weekday := range weekdays {
		if weekday == "WeekEnd" {
			// WeekEnd: Saturday (6) and Sunday (0)
			if today == time.Saturday || today == time.Sunday {
				return true
			}
		} else if weekday == "WeekDay" {
			// WeekDay: Monday (1) through Friday (5)
			if today >= time.Monday && today <= time.Friday {
				return true
			}
		}
	}

	return false
}
