package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"

	"github.com/cloudinn/escpos"
	"github.com/cloudinn/escpos/raster"
)

// ImagePrinter handles printing images to the printer
type ImagePrinter struct {
	printerName string
}

// NewImagePrinter creates a new image printer
func NewImagePrinter(printerName string) *ImagePrinter {
	return &ImagePrinter{
		printerName: printerName,
	}
}

// PrintImage prints an image using CloudInn/escpos library via lp command
func (ip *ImagePrinter) PrintImage(imagePath string, printerName string) error {
	// Check if the image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}

	// Use lp command with raw output for ESC/POS printing
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// Print header text
	stdin.Write([]byte("Printing image from: " + imagePath + "\n\n"))

	// Load the image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	// Decode the image (supports PNG, JPEG, GIF)
	img, imgFormat, err := image.Decode(imageFile)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to decode image: %v", err)
	}

	fmt.Printf("   📊 Loaded image, format: %s\n", imgFormat)
	fmt.Printf("   📊 Image dimensions: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())

	// Create a wrapper to make stdin compatible with io.ReadWriter
	readWriter := &readWriterWrapper{writer: stdin}

	// Create ESC/POS printer instance
	ep, err := escpos.NewPrinter(readWriter)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to create printer: %v", err)
	}

	// Initialize printer
	ep.Init()
	ep.SetAlign("center")

	// Create raster converter with Pi Zero 2W optimizations
	rasterConv := &raster.Converter{
		MaxWidth:  384, // Standard thermal printer width
		Threshold: 0.5, // Black/white threshold
	}

	fmt.Println("   🔄 Converting and printing image...")

	// Print the image using the raster converter
	rasterConv.Print(img, ep)

	// Add some spacing and cut
	ep.Linefeed()
	ep.Linefeed()
	ep.Write([]byte("Print job completed"))
	ep.Linefeed()
	ep.Cut()
	ep.End()

	// Close stdin and run the command
	stdin.Close()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("lp command failed: %v", err)
	}

	fmt.Println("   ✅ Print job completed successfully")
	return nil
}

// readWriterWrapper makes io.WriteCloser compatible with io.ReadWriter
type readWriterWrapper struct {
	writer io.WriteCloser
}

func (r *readWriterWrapper) Write(p []byte) (n int, err error) {
	return r.writer.Write(p)
}

func (r *readWriterWrapper) Read(p []byte) (n int, err error) {
	// For ESC/POS printing, we only need to write, not read
	return 0, io.EOF
}
