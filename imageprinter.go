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

	// DEBUG: Try different ESC/POS bitmap commands
	fmt.Println("   🔍 DEBUG: Testing different ESC/POS bitmap commands...")

	// Try ESC * m nL nH format (most common)
	fmt.Printf("   🔍 DEBUG: Trying ESC * 0 (8-dot single density)\n")
	escPosCmd := []byte("\x1B\x2A\x00") // ESC * 0
	fmt.Printf("   🔍 DEBUG: Command bytes: %02X %02X %02X\n", escPosCmd[0], escPosCmd[1], escPosCmd[2])
	stdin.Write(escPosCmd)

	// Send width (nL nH - low byte, high byte)
	widthBytes := []byte{byte(bytesPerLine & 0xFF), byte((bytesPerLine >> 8) & 0xFF)}
	fmt.Printf("   🔍 DEBUG: Width bytes: %02X %02X\n", widthBytes[0], widthBytes[1])
	stdin.Write(widthBytes)

	// Try alternative: ESC * m nL nH with different density
	fmt.Println("   🔍 DEBUG: Trying ESC * 1 (8-dot double density)")
	escPosCmd2 := []byte("\x1B\x2A\x01") // ESC * 1
	fmt.Printf("   🔍 DEBUG: Command bytes: %02X %02X %02X\n", escPosCmd2[0], escPosCmd2[1], escPosCmd2[2])
	stdin.Write(escPosCmd2)
	stdin.Write(widthBytes)

	// Try alternative: ESC * m nL nH with 24-dot density
	fmt.Println("   🔍 DEBUG: Trying ESC * 32 (24-dot single density)")
	escPosCmd3 := []byte("\x1B\x2A\x20") // ESC * 32 (0x20 = 32)
	fmt.Printf("   🔍 DEBUG: Command bytes: %02X %02X %02X\n", escPosCmd3[0], escPosCmd3[1], escPosCmd3[2])
	stdin.Write(escPosCmd3)
	stdin.Write(widthBytes)

	// Try alternative: GS v 0 command (some printers use this)
	fmt.Println("   🔍 DEBUG: Trying GS v 0 (alternative bitmap command)")
	gsCmd := []byte("\x1D\x76\x30\x00") // GS v 0
	fmt.Printf("   🔍 DEBUG: Command bytes: %02X %02X %02X %02X\n", gsCmd[0], gsCmd[1], gsCmd[2], gsCmd[3])
	stdin.Write(gsCmd)

	// Try alternative: GS * command (some printers use this)
	fmt.Println("   🔍 DEBUG: Trying GS * (alternative bitmap command)")
	gsCmd2 := []byte("\x1D\x2A\x00") // GS * 0
	fmt.Printf("   🔍 DEBUG: Command bytes: %02X %02X %02X\n", gsCmd2[0], gsCmd2[1], gsCmd2[2])
	stdin.Write(gsCmd2)
	stdin.Write(widthBytes)

	fmt.Println("   🔄 Processing image lines...")

	// DEBUG: Test with a simple pattern first
	fmt.Println("   🔍 DEBUG: Testing with simple pattern...")
	testPattern := make([]byte, bytesPerLine)
	// Create a simple test pattern: alternating black/white dots
	for i := 0; i < bytesPerLine; i++ {
		if i%2 == 0 {
			testPattern[i] = 0xAA // 10101010 pattern
		} else {
			testPattern[i] = 0x55 // 01010101 pattern
		}
	}
	fmt.Printf("   🔍 DEBUG: Test pattern (first 10 bytes): ")
	for i := 0; i < 10 && i < len(testPattern); i++ {
		fmt.Printf("%02X ", testPattern[i])
	}
	fmt.Println()

	// Send test pattern
	stdin.Write(testPattern)

	// DEBUG: Sample first few pixels to check conversion
	fmt.Println("   🔍 DEBUG: Sampling first 10 pixels:")
	for x := 0; x < 10 && x < width; x++ {
		r, g, b, _ := img.At(x, 0).RGBA()
		isBlack := r < 32768 || g < 32768 || b < 32768
		fmt.Printf("   🔍   Pixel %d: RGB(%d,%d,%d) -> %s\n", x, r>>8, g>>8, b>>8, map[bool]string{true: "BLACK", false: "WHITE"}[isBlack])
	}

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

		// DEBUG: Show first line's bitmap data
		if y == 0 {
			fmt.Printf("   🔍 DEBUG: First line bitmap data (first 10 bytes): ")
			for i := 0; i < 10 && i < len(lineData); i++ {
				fmt.Printf("%02X ", lineData[i])
			}
			fmt.Println()
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
