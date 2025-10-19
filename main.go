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

	// File path for PBM
	pbmPath := filepath.Join(outputDir, "output.pbm")

	fmt.Println("🔄 Converting SVG to PBM...")
	if err := ConvertSVGToImage(pbmPath); err != nil {
		fmt.Printf("Error converting SVG to PBM: %v\n", err)
		return
	}
	fmt.Println("✅ SVG to PBM conversion completed")

	fmt.Println("🔄 Printing PBM...")
	if err := imagePrinter.PrintImage(pbmPath); err != nil {
		fmt.Printf("Error printing PBM: %v\n", err)
		return
	}
	fmt.Println("✅ Print job sent successfully!")

	// Clean up temporary files
	defer func() {
		os.Remove(pbmPath)
	}()
}
