package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ConvertSVGToImage converts SVG template to JPEG image using ImageMagick
func ConvertSVGToImage(outputPath string) error {
	// Check if ImageMagick is available
	if _, err := exec.LookPath("convert"); err != nil {
		return fmt.Errorf("ImageMagick not found. Please install it with: sudo apt-get install imagemagick")
	}

	// Load SVG content from templates folder
	svgContent, err := loadSVGFromTemplates()
	if err != nil {
		return fmt.Errorf("failed to load SVG from templates: %v", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create temporary SVG file
	tempSVGPath := filepath.Join(outputDir, "temp.svg")
	if err := os.WriteFile(tempSVGPath, []byte(svgContent), 0644); err != nil {
		return fmt.Errorf("failed to write temporary SVG file: %v", err)
	}
	defer os.Remove(tempSVGPath) // Clean up temp file

	// Convert SVG to JPEG using ImageMagick with Pi optimizations
	// Force width to 384 dots (standard thermal printer width) and maintain aspect ratio
	cmd := exec.Command("convert",
		"-background", "white",
		"-density", "150", // Reduced density for Pi
		"-resize", "384x", // Force width to 384 dots, height auto-calculated
		"-colorspace", "Gray", // Convert to grayscale for faster processing
		"-threshold", "50%", // Convert to black/white for thermal printer
		"-quality", "75", // Lower quality for faster processing
		"-limit", "memory", "64MB", // Even more aggressive memory limits for Pi Zero
		"-limit", "map", "128MB", // Reduced memory mapping
		"-limit", "disk", "256MB", // Reduced disk usage
		"-define", "registry:temporary-path=/tmp", // Use /tmp for temp files
		tempSVGPath,
		outputPath,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ImageMagick conversion failed: %v", err)
	}

	return nil
}

// loadSVGFromTemplates loads SVG content from the templates folder
func loadSVGFromTemplates() (string, error) {
	// Get template path relative to executable directory
	templatePath, err := GetExecutableRelativePath("templates/sample.svg")
	if err != nil {
		return "", fmt.Errorf("failed to get template path: %v", err)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %v", templatePath, err)
	}

	// Replace timestamp placeholder
	svgContent := string(content)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	svgContent = strings.ReplaceAll(svgContent, "{{.Timestamp}}", timestamp)

	return svgContent, nil
}
