package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Configuration
	printerName := os.Getenv("PRINTER_NAME")

	// Get output directory relative to executable
	outputDir, err := GetExecutableRelativePath("tmp/printy")
	if err != nil {
		fmt.Printf("Error getting output directory: %v\n", err)
		return
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Initialize printer
	imagePrinter := NewImagePrinter(printerName)

	// File path
	imagePath := filepath.Join(outputDir, "output.png")

	fmt.Println("🔄 Converting HTML to image...")
	if err := ConvertHTMLToImage(imagePath); err != nil {
		fmt.Printf("Error converting HTML to image: %v\n", err)
		return
	}
	fmt.Println("✅ HTML to image conversion completed")

	fmt.Println("🔄 Printing image...")
	if err := imagePrinter.PrintImage(imagePath); err != nil {
		fmt.Printf("Error printing image: %v\n", err)
		return
	}
	fmt.Println("✅ Print job sent successfully!")

	// Clean up Chrome instance when done
	defer CloseChromeInstance()

	// Clean up temporary files
	defer func() {
		os.Remove(imagePath)
	}()
}
