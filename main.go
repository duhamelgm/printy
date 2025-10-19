package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Configuration
	printerName := os.Getenv("PRINTER_NAME")
	if printerName == "" {
		fmt.Println("⚠️  PRINTER_NAME environment variable not set, will use default printer")
	} else {
		fmt.Printf("🖨️  Printer configured: %s\n", printerName)
	}

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

	// File path for PNG (faster conversion with rsvg-convert)
	pngPath := filepath.Join(outputDir, "output.png")

	fmt.Println("🔄 Converting SVG to PNG...")
	if err := ConvertSVGToImage(pngPath); err != nil {
		fmt.Printf("Error converting SVG to PNG: %v\n", err)
		return
	}
	fmt.Println("✅ SVG to PNG conversion completed")

	fmt.Println("🔄 Printing with ESC/POS...")
	if err := imagePrinter.PrintImage(pngPath, printerName); err != nil {
		fmt.Printf("Error printing with ESC/POS: %v\n", err)
		return
	}
	fmt.Println("✅ ESC/POS print job sent successfully!")

	// Keep PNG file for testing
	fmt.Printf("📁 PNG file saved at: %s\n", pngPath)
	fmt.Println("🔧 You can test ESC/POS printing manually:")
	fmt.Println("   Test direct device access:")
	fmt.Printf("   echo 'Hello World' > /dev/usb/lp0\n")
	fmt.Println("   Test with different device paths:")
	fmt.Println("   /dev/usb/lp0, /dev/usb/lp1, /dev/lp0, /dev/lp1")
}
