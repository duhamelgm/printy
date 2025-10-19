package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

// Image conversion options - defined in the file
const (
	ImageQuality = 90
	ImageFormat  = "png"
	ImageTimeout = 30 * time.Second
)

// Global Chrome instance for reuse
var (
	chromeCtx    context.Context
	chromeCancel context.CancelFunc
	chromeMutex  sync.Mutex
	chromeReady  bool
)

// getOrCreateChromeInstance gets the existing Chrome instance or creates a new one
func getOrCreateChromeInstance() (context.Context, error) {
	chromeMutex.Lock()
	defer chromeMutex.Unlock()

	if chromeReady && chromeCtx != nil {
		// Check if context is still valid
		select {
		case <-chromeCtx.Done():
			// Context is cancelled, need to recreate
			chromeReady = false
		default:
			// Context is still valid, reuse it
			return chromeCtx, nil
		}
	}

	// Create new Chrome instance
	ctx, _ := context.WithTimeout(context.Background(), ImageTimeout)

	// Get Chrome data directory relative to executable
	chromeDataDir, err := GetExecutableRelativePath("tmp/chrome-data")
	if err != nil {
		return nil, fmt.Errorf("failed to get Chrome data directory: %v", err)
	}

	// Create ChromeDP context with unique user data dir
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.UserDataDir(chromeDataDir), // Unique data directory
	)

	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	chromeCtx, chromeCancel = chromedp.NewContext(allocCtx)

	chromeReady = true
	return chromeCtx, nil
}

// CloseChromeInstance closes the global Chrome instance
func CloseChromeInstance() {
	chromeMutex.Lock()
	defer chromeMutex.Unlock()

	if chromeCancel != nil {
		chromeCancel()
		chromeCancel = nil
	}
	chromeCtx = nil
	chromeReady = false
}

// ConvertHTMLToImage converts HTML from templates folder to an image using ChromeDP
func ConvertHTMLToImage(outputPath string) error {
	// Load HTML content from templates folder
	htmlContent, err := loadHTMLFromTemplates()
	if err != nil {
		return fmt.Errorf("failed to load HTML from templates: %v", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Get or create Chrome instance
	ctx, err := getOrCreateChromeInstance()
	if err != nil {
		return fmt.Errorf("failed to get Chrome instance: %v", err)
	}

	// Capture full page screenshot
	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+htmlContent),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(1*time.Second), // Wait for any dynamic content
		chromedp.FullScreenshot(&buf, ImageQuality),
	)

	if err != nil {
		return fmt.Errorf("ChromeDP screenshot failed: %v", err)
	}

	// Write the image to file
	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %v", err)
	}

	return nil
}

// loadHTMLFromTemplates loads HTML content from the templates folder
func loadHTMLFromTemplates() (string, error) {
	// Get template path relative to executable directory
	templatePath, err := GetExecutableRelativePath("templates/sample.html")
	if err != nil {
		return "", fmt.Errorf("failed to get template path: %v", err)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %v", templatePath, err)
	}

	return string(content), nil
}
