package printer

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

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

	// Load the image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	// Decode the image (supports PNG, JPEG, GIF)
	img, _, err := image.Decode(imageFile)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to decode image: %v", err)
	}

	img = invertColors(img)

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
	ep.SetAlign("left")

	// Create raster converter with Pi Zero 2W optimizations
	rasterConv := &raster.Converter{
		MaxWidth:  576, // Standard thermal printer width
		Threshold: 0.5, // Black/white threshold
	}

	sz := img.Bounds().Size()

	data, rw, bw := rasterConv.ToRaster(img)

	ep.Raster(rw, sz.Y, bw, data, "bitImage")

	// Print the image using the raster converter
	// rasterConv.Print(img, ep)

	// Add some spacing and cut
	ep.Linefeed()
	ep.Cut()
	ep.End()

	// Close stdin and run the command
	stdin.Close()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("lp command failed: %v", err)
	}

	return nil
}

func intLowHigh(inpNumber int, outBytes int) []byte {

	maxInput := (256 << (uint((outBytes * 8)) - 1))

	if outBytes < 1 || outBytes > 4 {
		log.Println("Can only output 1-4 bytes")
	}
	if inpNumber < 0 || inpNumber > maxInput {
		log.Printf("Number too large. Can only output up to %d in %d bytes\n", maxInput, outBytes)
	}
	var outp []byte
	for i := 0; i < outBytes; i++ {
		inpNumberByte := byte(inpNumber % 256)
		outp = append(outp, inpNumberByte)
		inpNumber = inpNumber / 256
	}
	return outp
}

// PrintRawRaster sends precomputed raster/ESC-POS binary data directly to the printer.
// The data is written as-is to the printer in raw mode.
func (ip *ImagePrinter) PrintRawRaster(rasterData []byte, width int, height int) error {
	if len(rasterData) == 0 {
		return fmt.Errorf("raster data is empty")
	}

	bytesPerRow := (width + 7) >> 3
	if bytesPerRow <= 0 {
		return fmt.Errorf("invalid raster width: %d", width)
	}

	expectedLen := bytesPerRow * height
	if len(rasterData) < expectedLen {
		return fmt.Errorf("raster data too short: have %d, need %d", len(rasterData), expectedLen)
	}

	densityByte := byte(0)
	headerBase := []byte{0x1D, 0x76, 0x30, densityByte}
	headerBase = append(headerBase, intLowHigh(bytesPerRow, 2)...)

	totalSegments := (height + 599) / 600
	log.Printf("Printing raster in %d segment(s)", totalSegments)

	cmd := exec.Command("lp", "-d", ip.printerName, "-o", "raw")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to start lp command: %v", err)
	}

	readWriter := &readWriterWrapper{writer: stdin}
	ep, err := escpos.NewPrinter(readWriter)
	if err != nil {
		stdin.Close()
		return fmt.Errorf("failed to create escpos printer: %v", err)
	}

	ep.Init()
	ep.SetAlign("left")

	for segment := 0; segment < totalSegments; segment++ {
		segmentStartRow := segment * 600
		segmentHeight := 600
		if segmentStartRow+segmentHeight > height {
			segmentHeight = height - segmentStartRow
		}

		startIdx := segmentStartRow * bytesPerRow
		endIdx := startIdx + segmentHeight*bytesPerRow
		if endIdx > len(rasterData) {
			stdin.Close()
			return fmt.Errorf("segment %d exceeds raster data length", segment)
		}

		header := append([]byte{}, headerBase...)
		header = append(header, intLowHigh(segmentHeight, 2)...)
		fullSegment := append(header, rasterData[startIdx:endIdx]...)

		if _, err := ep.Write(fullSegment); err != nil {
			stdin.Close()
			return fmt.Errorf("failed to write raster data for segment %d: %v", segment, err)
		}

		if segment < totalSegments-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	ep.Linefeed()
	ep.Cut()
	ep.End()

	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("lp command failed: %v", err)
	}

	return nil
}

// invertColors inverts the colors of an image for thermal printing
func invertColors(img image.Image) image.Image {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// Invert RGB values (255 - value)
			invertedColor := color.RGBA{
				R: 255 - uint8(r>>8),
				G: 255 - uint8(g>>8),
				B: 255 - uint8(b>>8),
				A: uint8(a >> 8),
			}
			inverted.Set(x, y, invertedColor)
		}
	}

	return inverted
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
