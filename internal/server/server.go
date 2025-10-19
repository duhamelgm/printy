package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"printy/internal/db"
	"printy/internal/notion"
	"printy/internal/printer"
	"printy/internal/tickets"
)

// Server represents the HTTP server
type Server struct {
	printer  *printer.Printer
	database *db.Database
	port     string
}

// PrintRequest represents the request body for printing
type PrintRequest struct {
	PrinterName string `json:"printer_name,omitempty"`
}

// PrintResponse represents the response for printing
type PrintResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SyncTicketsResponse represents the response for syncing tickets
type SyncTicketsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count"`
	Error   string `json:"error,omitempty"`
}

// New creates a new HTTP server
func New(port string) (*Server, error) {
	printerName := os.Getenv("PRINTER_NAME")

	// Initialize printer
	p, err := printer.New(printerName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize printer: %v", err)
	}

	execDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	database, err := db.New(filepath.Join(execDir, "data", "printy.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return &Server{
		printer:  p,
		database: database,
		port:     port,
	}, nil
}

// Close closes the database connection
func (s *Server) Close() error {
	if s.database != nil {
		return s.database.Close()
	}
	return nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/print", s.handlePrint)
	mux.HandleFunc("/print/", s.handlePrint) // Handle trailing slash
	mux.HandleFunc("/sync-tickets", s.handleSyncTickets)
	mux.HandleFunc("/sync-tickets/", s.handleSyncTickets) // Handle trailing slash

	server := &http.Server{
		Addr:         ":" + s.port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}

// handlePrint handles print requests
func (s *Server) handlePrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get relevant tickets for today
	relevantTickets, err := tickets.GetRelevantTickets(s.database, 10)
	if err != nil {
		response := PrintResponse{
			Success: false,
			Message: "Failed to get relevant tickets",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Print each relevant ticket
	for _, ticket := range relevantTickets {
		if err := s.printer.Print(ticket.RefID, ticket.Title); err != nil {
			log.Printf("Failed to print ticket %d: %v", ticket.ID, err)
			continue
		}

		// Create a print record in database
		print := &db.Print{
			TicketID:  ticket.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.database.CreatePrint(print); err != nil {
			log.Printf("⚠️  Warning: Failed to record print for ticket %d: %v", ticket.ID, err)
		}
	}

	// Success response
	response := PrintResponse{
		Success: true,
		Message: fmt.Sprintf("Print job completed successfully for %d tickets", len(relevantTickets)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleSyncTickets handles sync tickets requests
func (s *Server) handleSyncTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get Notion API credentials from environment
	apiKey := os.Getenv("NOTION_API_KEY")
	databaseID := os.Getenv("NOTION_DATABASE_ID")

	if apiKey == "" || databaseID == "" {
		response := SyncTicketsResponse{
			Success: false,
			Message: "Notion API credentials not configured",
			Error:   "NOTION_API_KEY and NOTION_DATABASE_ID environment variables must be set",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Fetch tickets from Notion
	notionTickets, err := notion.GetTickets(apiKey, databaseID)
	if err != nil {
		response := SyncTicketsResponse{
			Success: false,
			Message: "Failed to fetch tickets from Notion",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Process and store tickets
	syncedCount := 0
	for _, notionTicket := range notionTickets {
		// Parse priority and cooldown using the parser module
		priority := tickets.ParsePriority(notionTicket.Priority)
		cooldown := tickets.ParseCooldown(notionTicket.Cooldown)
		weekdays := notionTicket.Weekdays // Keep as JSON array string

		// Check if ticket already exists
		existingTicket, err := s.database.GetTicketByRefID(notionTicket.ID)
		if err != nil && err.Error() != "ticket not found" {
			log.Printf("Error checking existing ticket %s: %v", notionTicket.ID, err)
			continue
		}

		now := time.Now()
		if existingTicket != nil {
			// Update existing ticket
			existingTicket.Title = notionTicket.Name
			existingTicket.Priority = priority
			existingTicket.Cooldown = cooldown
			existingTicket.Weekdays = weekdays
			existingTicket.UpdatedAt = now

			if err := s.database.UpdateTicket(existingTicket); err != nil {
				log.Printf("Error updating ticket %s: %v", notionTicket.ID, err)
				continue
			}
		} else {
			// Create new ticket
			ticket := &db.Ticket{
				RefID:     notionTicket.ID,
				Title:     notionTicket.Name,
				Priority:  priority,
				Cooldown:  cooldown,
				Weekdays:  weekdays,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := s.database.CreateTicket(ticket); err != nil {
				log.Printf("Error creating ticket %s: %v", notionTicket.ID, err)
				continue
			}
		}

		syncedCount++
	}

	// Success response
	response := SyncTicketsResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully synced %d tickets from Notion", syncedCount),
		Count:   syncedCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
