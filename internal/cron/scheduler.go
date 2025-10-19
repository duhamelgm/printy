package cron

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Scheduler handles cron-like scheduling for print jobs
type Scheduler struct {
	client *http.Client
	url    string
}

// NewScheduler creates a new scheduler instance
func NewScheduler(baseURL string) *Scheduler {
	return &Scheduler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		url: baseURL,
	}
}

// StartDailyPrintJob starts a goroutine that calls print-backlog daily at 8:00 AM Montreal time
func (s *Scheduler) StartDailyPrintJob() {
	go func() {
		// Set timezone to Montreal
		location, err := time.LoadLocation("America/Montreal")
		if err != nil {
			log.Printf("Failed to load Montreal timezone: %v", err)
			return
		}

		for {
			now := time.Now().In(location)

			// Calculate next 8:00 AM Montreal time
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, location)
			if now.After(nextRun) {
				// If it's already past 8:00 AM today, schedule for tomorrow
				nextRun = nextRun.Add(24 * time.Hour)
			}

			// Wait until next run time
			waitDuration := time.Until(nextRun)
			log.Printf("Next print job scheduled for: %s (in %v)", nextRun.Format("2006-01-02 15:04:05 MST"), waitDuration)

			time.Sleep(waitDuration)

			// Execute the print job
			s.executePrintJob()
		}
	}()
}

// executePrintJob calls the print-backlog endpoint
func (s *Scheduler) executePrintJob() {
	log.Printf("Executing daily print job at %s", time.Now().Format("2006-01-02 15:04:05 MST"))

	// Create JSON request body
	requestBody := map[string]interface{}{
		"assignee": "Duhamel",
		"count":    3,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("ERROR: Failed to marshal JSON: %v", err)
		return
	}

	resp, err := s.client.Post(s.url+"/print-backlog", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ERROR: Failed to call print-backlog endpoint: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Print job completed successfully (status: %d)", resp.StatusCode)
	} else {
		log.Printf("Print job failed with status: %d", resp.StatusCode)
	}
}
