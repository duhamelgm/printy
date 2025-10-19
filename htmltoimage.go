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
	ImageQuality      = 90
	ImageFormat       = "png"
	ImageTimeout      = 30 * time.Second
	BrowserTimeout    = 90 * time.Second // Increased timeout for slower Pi hardware
	ScreenshotTimeout = 60 * time.Second // Separate timeout for screenshot operations on Pi
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

	// Create new Chrome instance with increased timeout
	ctx, cancel := context.WithTimeout(context.Background(), BrowserTimeout)
	defer cancel()

	// Get Chrome data directory relative to executable
	chromeDataDir, err := GetExecutableRelativePath("tmp/chrome-data")
	if err != nil {
		return nil, fmt.Errorf("failed to get Chrome data directory: %v", err)
	}

	// Create ChromeDP context optimized for Raspberry Pi with 512MB RAM
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		// Keep CSS enabled but disable JavaScript for Pi performance
		chromedp.Flag("disable-javascript", true),                                 // Keep JavaScript disabled for performance
		chromedp.Flag("disable-plugins", true),                                    // Disable plugins
		chromedp.Flag("disable-extensions", true),                                 // Disable extensions
		chromedp.Flag("disable-background-timer-throttling", true),                // Disable background throttling
		chromedp.Flag("disable-backgrounding-occluded-windows", true),             // Disable backgrounding
		chromedp.Flag("disable-renderer-backgrounding", true),                     // Disable renderer backgrounding
		chromedp.Flag("disable-background-networking", true),                      // Disable background networking
		chromedp.Flag("disable-default-apps", true),                               // Disable default apps
		chromedp.Flag("disable-sync", true),                                       // Disable sync
		chromedp.Flag("disable-translate", true),                                  // Disable translate
		chromedp.Flag("disable-ipc-flooding-protection", true),                    // Disable IPC flooding protection
		chromedp.Flag("disable-hang-monitor", true),                               // Disable hang monitor
		chromedp.Flag("disable-prompt-on-repost", true),                           // Disable prompt on repost
		chromedp.Flag("disable-domain-reliability", true),                         // Disable domain reliability
		chromedp.Flag("disable-component-extensions-with-background-pages", true), // Disable component extensions
		chromedp.Flag("aggressive-cache-discard", true),                           // Aggressive cache discard
		chromedp.Flag("memory-pressure-off", true),                                // Turn off memory pressure
		// Pi-specific memory optimizations
		chromedp.Flag("max_old_space_size", "256"),                    // Reduced memory limit for Pi (256MB)
		chromedp.Flag("single-process", true),                         // Single process mode for Pi
		chromedp.Flag("disable-software-rasterizer", true),            // Disable software rasterizer
		chromedp.Flag("disable-threaded-compositing", true),           // Disable threaded compositing
		chromedp.Flag("disable-threaded-animation", true),             // Disable threaded animation
		chromedp.Flag("disable-checker-imaging", true),                // Disable checker imaging
		chromedp.Flag("disable-new-tab-first-run", true),              // Disable new tab first run
		chromedp.Flag("disable-client-side-phishing-detection", true), // Disable phishing detection
		chromedp.Flag("disable-component-update", true),               // Disable component updates
		chromedp.Flag("disable-background-mode", true),                // Disable background mode
		chromedp.Flag("disable-logging", true),                        // Disable logging
		chromedp.Flag("silent", true),                                 // Silent mode
		chromedp.Flag("disable-default-browser-check", true),          // Disable default browser check
		chromedp.UserDataDir(chromeDataDir),                           // Unique data directory
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

	// Create a context with timeout specifically for screenshot operations
	screenshotCtx, cancel := context.WithTimeout(ctx, ScreenshotTimeout)
	defer cancel()

	// Capture full page screenshot with optimized settings
	var buf []byte
	err = chromedp.Run(screenshotCtx,
		chromedp.Navigate("data:text/html,"+htmlContent),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second), // Increased wait time for Pi's slower hardware
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
