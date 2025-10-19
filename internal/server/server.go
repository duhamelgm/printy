package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"printy/internal/printer"
)

// Server represents the HTTP server
type Server struct {
	printer *printer.Printer
	port    string
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

// New creates a new HTTP server
func New(port string) (*Server, error) {
	// Get printer name from environment or use default
	printerName := os.Getenv("PRINTER_NAME")
	if printerName == "" {
		printerName = "default"
	}

	// Initialize printer
	p, err := printer.New(printerName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize printer: %v", err)
	}

	return &Server{
		printer: p,
		port:    port,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/print", s.handlePrint)
	mux.HandleFunc("/print/", s.handlePrint) // Handle trailing slash

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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.printer.Print(); err != nil {
		response := PrintResponse{
			Success: false,
			Message: "Print job failed",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := PrintResponse{
		Success: true,
		Message: "Print job completed successfully for printer",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
