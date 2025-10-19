package printer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ConvertSVGToImage converts SVG template to JPEG image using rsvg-convert
func ConvertSVGToImage(outputPath string) error {
	// Check if rsvg-convert is available
	if _, err := exec.LookPath("rsvg-convert"); err != nil {
		return fmt.Errorf("rsvg-convert not found. Please install it with: sudo apt-get install librsvg2-bin")
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

	// Convert SVG to PNG using rsvg-convert (much faster than ImageMagick)
	// Force width to 384 dots (standard thermal printer width) and maintain aspect ratio
	cmd := exec.Command("rsvg-convert",
		"--width", "512", // Force width to 384 dots, height auto-calculated
		"--format", "png", // Output PNG format directly
		"--output", outputPath,
		tempSVGPath,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsvg-convert SVG to PNG conversion failed: %v", err)
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
