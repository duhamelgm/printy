package main

import (
	"bytes"
	"fmt"
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

// PrintImage prints an image file to the specified printer
func (ip *ImagePrinter) PrintImage(imagePath string) error {
	// Check if the image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}

	// Print using CUPS with lp command
	cmd := exec.Command("lp", "-d", ip.printerName, imagePath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to print image: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
