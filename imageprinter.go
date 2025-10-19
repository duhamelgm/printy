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

// PrintImage prints a PBM file to the specified printer in raw mode
func (ip *ImagePrinter) PrintImage(pbmPath string) error {
	// Check if the PBM file exists
	if _, err := os.Stat(pbmPath); os.IsNotExist(err) {
		return fmt.Errorf("PBM file does not exist: %s", pbmPath)
	}

	// Send the PBM file directly to the printer using the printer name
	// This assumes the printer is set up in CUPS and can handle raw data
	cmd := exec.Command("lp", "-d", ip.printerName, "-o", "raw", pbmPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to print PBM file: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
