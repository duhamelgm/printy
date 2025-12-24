package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"printy/internal/queue"
	"printy/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Parse command line flags
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// Create and start server
	s, err := server.New(*port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer s.Close() // Ensure database connection is closed

	// Start queue runner
	queueRunner, runnerErr := queue.NewRunner()
	if runnerErr != nil {
		log.Printf("Queue runner not started: %v", runnerErr)
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go queueRunner.Run(ctx)
		defer queueRunner.Close()
		log.Printf("Queue runner listening on Redis queue %s", "printing_queue")
	}

	fmt.Printf("🚀 Starting Printy HTTP Server on port %s\n", *port)
	fmt.Printf("📋 Set PRINTER_NAME environment variable to specify printer\n")
	fmt.Printf("📊 Set DB_PATH environment variable to specify database location\n")
	fmt.Printf("🌐 Server will be available at: http://localhost:%s\n", *port)

	// Start the daily print job scheduler
	// scheduler := cron.NewScheduler(fmt.Sprintf("http://localhost:%s", *port))
	// scheduler.StartDailyPrintJob()
	// fmt.Printf("⏰ Daily print job scheduled for 8:00 AM Montreal time\n")

	// Start server (this blocks) so the queue runner stays alive
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
