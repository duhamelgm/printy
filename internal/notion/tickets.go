package notion

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TicketItem represents a formatted ticket item
type TicketItem struct {
	Cooldown string
	ID       string
	Priority string
	Weekdays string
	Name     string
	Assignee string
}

// GetTickets fetches items from Notion database and formats them
func GetTickets(apiKey, databaseID string) ([]TicketItem, error) {
	client := NewClient(apiKey)

	pages, err := client.QueryDatabase(databaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %v", err)
	}

	var tickets []TicketItem

	for _, page := range pages {
		properties, ok := page["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		ticket := TicketItem{}

		// Extract cooldown (rich_text)
		if cooldownProp, ok := properties["cooldown"].(map[string]interface{}); ok {
			if richText, ok := cooldownProp["rich_text"].([]interface{}); ok && len(richText) > 0 {
				if textObj, ok := richText[0].(map[string]interface{}); ok {
					if plainText, ok := textObj["plain_text"].(string); ok {
						ticket.Cooldown = plainText
					}
				}
			}
		}

		// Extract id (unique_id)
		if idProp, ok := properties["id"].(map[string]interface{}); ok {
			if uniqueID, ok := idProp["unique_id"].(map[string]interface{}); ok {
				if number, ok := uniqueID["number"].(float64); ok {
					// Get the prefix from Notion's unique_id
					if prefix, ok := uniqueID["prefix"].(string); ok {
						ticket.ID = fmt.Sprintf("%s%.0f", prefix, number)
					} else {
						// Fallback if no prefix is provided
						ticket.ID = fmt.Sprintf("%.0f", number)
					}
				}
			}
		}

		// Extract priority (select)
		if priorityProp, ok := properties["priority"].(map[string]interface{}); ok {
			if selectObj, ok := priorityProp["select"].(map[string]interface{}); ok {
				if name, ok := selectObj["name"].(string); ok {
					ticket.Priority = name
				}
			}
		}

		// Extract weekdays (multi_select)
		if weekdaysProp, ok := properties["weekdays"].(map[string]interface{}); ok {
			if multiSelect, ok := weekdaysProp["multi_select"].([]interface{}); ok {
				var weekdayNames []string
				for _, item := range multiSelect {
					if itemObj, ok := item.(map[string]interface{}); ok {
						if name, ok := itemObj["name"].(string); ok {
							// Only include valid weekday types
							if name == "WeekEnd" || name == "WeekDay" {
								weekdayNames = append(weekdayNames, name)
							}
						}
					}
				}
				// Store as JSON array string
				if len(weekdayNames) > 0 {
					jsonData, err := json.Marshal(weekdayNames)
					if err != nil {
						// Fallback to comma-separated if JSON marshaling fails
						ticket.Weekdays = strings.Join(weekdayNames, ", ")
					} else {
						ticket.Weekdays = string(jsonData)
					}
				} else {
					ticket.Weekdays = "[]"
				}
			}
		}

		// Extract name (title)
		if nameProp, ok := properties["name"].(map[string]interface{}); ok {
			if title, ok := nameProp["title"].([]interface{}); ok && len(title) > 0 {
				if titleObj, ok := title[0].(map[string]interface{}); ok {
					if plainText, ok := titleObj["plain_text"].(string); ok {
						ticket.Name = plainText
					}
				}
			}
		}

		// Extract assignee (people) - save first names of all assignees
		if assigneeProp, ok := properties["assignee"].(map[string]interface{}); ok {
			if people, ok := assigneeProp["people"].([]interface{}); ok && len(people) > 0 {
				var firstNames []string
				for _, personInterface := range people {
					if person, ok := personInterface.(map[string]interface{}); ok {
						if name, ok := person["name"].(string); ok {
							// Extract only the first name
							parts := strings.Fields(name)
							if len(parts) > 0 {
								firstNames = append(firstNames, parts[0])
							}
						}
					}
				}
				// Join all first names with commas
				ticket.Assignee = strings.Join(firstNames, ", ")
			}
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
