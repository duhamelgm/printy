package main

import (
	"flag"
	"fmt"
	"log"

	"printy/internal/server"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// Create and start server
	s, err := server.New(*port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Printf("🚀 Starting Printy HTTP Server on port %s\n", *port)

	// Start server (this blocks)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
