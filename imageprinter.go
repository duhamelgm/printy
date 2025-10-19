package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
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

// PrintImage prints a JPEG image using ESC/POS commands via lp command
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

	// Print header text using basic ESC/POS commands
	stdin.Write([]byte("Printing image from: " + imagePath + "\n\n"))

	// Load the JPEG image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	// Decode the JPEG image (much faster than PNG on Pi Zero 2W)
	img, err := jpeg.Decode(imageFile)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to decode JPEG image: %v", err)
	}

	// Print the image using basic ESC/POS bitmap commands
	if err := printImageToPrinter(stdin, img); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to print image: %v", err)
	}

	// Add some spacing and cut
	stdin.Write([]byte("\n\nPrint job completed\n"))
	stdin.Write([]byte("\x1D\x56\x00")) // ESC/POS cut command

	// Close stdin and run the command
	stdin.Close()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("lp command failed: %v", err)
	}

	return nil
}

// printImageToPrinter converts an image to ESC/POS bitmap format and sends it to the printer
func printImageToPrinter(stdin io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("   📊 Image dimensions: %dx%d pixels\n", width, height)

	// ESC/POS bitmap command: ESC * m nL nH d1...dk
	// m = 0 (8-dot single density), 1 (8-dot double density), 32 (24-dot single density), 33 (24-dot double density)
	// For simplicity, we'll use 8-dot single density (m=0)

	// Calculate bytes per line (8 dots per byte)
	bytesPerLine := (width + 7) / 8
	fmt.Printf("   📊 Bytes per line: %d\n", bytesPerLine)
	fmt.Printf("   📊 Total lines to process: %d\n", height)

	// Send ESC/POS bitmap command
	stdin.Write([]byte("\x1B\x2A\x00")) // ESC * 0 (8-dot single density)

	// Send width (nL nH - low byte, high byte)
	stdin.Write([]byte{byte(bytesPerLine & 0xFF), byte((bytesPerLine >> 8) & 0xFF)})

	fmt.Println("   🔄 Processing image lines...")

	// Convert image to bitmap data (optimized for Pi Zero 2W)
	// Since ImageMagick already converted to black/white, we can optimize further
	for y := 0; y < height; y++ {
		lineData := make([]byte, bytesPerLine)

		for x := 0; x < width; x++ {
			// Get pixel color (already black/white from ImageMagick)
			r, g, b, _ := img.At(x, y).RGBA()

			// Fast black/white check (since ImageMagick already thresholded)
			// Just check if any color component is below threshold
			isBlack := r < 32768 || g < 32768 || b < 32768

			if isBlack {
				// Set corresponding bit in the byte
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				lineData[byteIndex] |= 1 << bitIndex
			}
		}

		// Send line data
		stdin.Write(lineData)

		// Progress indicator for every 10% of lines
		if (y+1)%(height/10+1) == 0 || y == height-1 {
			progress := ((y + 1) * 100) / height
			fmt.Printf("   📈 Progress: %d%% (%d/%d lines)\n", progress, y+1, height)
		}
	}

	fmt.Println("   ✅ Bitmap conversion completed")
	return nil
}
