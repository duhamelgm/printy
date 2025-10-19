package printer

import (
	"fmt"
	"os"
	"path/filepath"
)

// Printer handles all printing operations
type Printer struct {
	printerName string
	outputDir   string
}

// New creates a new printer instance
func New(printerName string) (*Printer, error) {
	// Get output directory relative to executable
	outputDir, err := GetExecutableRelativePath("tmp/printy")
	if err != nil {
		return nil, fmt.Errorf("error getting output directory: %v", err)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory: %v", err)
	}

	return &Printer{
		printerName: printerName,
		outputDir:   outputDir,
	}, nil
}

// Print executes the complete printing workflow
func (p *Printer) Print(ticketID, title string) error {
	// File path for PNG (faster conversion with rsvg-convert)
	pngPath := filepath.Join(p.outputDir, "output.png")

	fmt.Println("🔄 Converting SVG to PNG...")
	if err := ConvertSVGToImage(pngPath, ticketID, title); err != nil {
		return fmt.Errorf("error converting SVG to PNG: %v", err)
	}

	imagePrinter := NewImagePrinter(p.printerName)
	if err := imagePrinter.PrintImage(pngPath, p.printerName); err != nil {
		return fmt.Errorf("error printing with ESC/POS: %v", err)
	}

	return nil
}
