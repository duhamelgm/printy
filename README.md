# Printy - HTML to Image Printing System

A Go-based system for converting HTML documents directly to images and printing them using ChromeDP. This system provides a streamlined approach to HTML-to-image conversion without intermediate PDF generation.

## Features

- 🚀 **Direct HTML to Image**: No intermediate PDF conversion
- ⚡ **Chrome Instance Reuse**: Fast subsequent conversions
- 🖨️ **CUPS Integration**: Direct printing support
- 📁 **Template-based**: HTML templates from files
- 🎯 **Full Page Capture**: Captures entire HTML content
- 🔧 **Pure Go**: No Node.js or external dependencies

## Prerequisites

Before running this application, you need to install the following dependencies:

### macOS
```bash
# Install Chrome/Chromium (usually already installed)
# If not installed:
brew install --cask google-chrome
```

### Ubuntu/Debian
```bash
# Install Chrome/Chromium
sudo apt-get update
sudo apt-get install chromium-browser
```

### Other Systems
- **Chrome/Chromium**: Download from [https://www.google.com/chrome/](https://www.google.com/chrome/)

## Project Structure

```
printy/
├── main.go              # Main application entry point
├── htmltoimage.go        # HTML to image conversion using ChromeDP
├── imageprinter.go        # Image printing module using CUPS
├── templates/
│   └── sample.html        # HTML template with Tailwind CSS
├── tmp/                   # Temporary files directory
│   ├── chrome-data/       # Chrome user data directory
│   └── printy/            # Output files
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd printy
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run .
   ```

## Usage

The application will:
1. Load HTML content from `templates/sample.html`
2. Convert HTML directly to image using ChromeDP
3. Print the image to your configured printer
4. Clean up temporary files

## Configuration

### Printer Configuration
The application is configured to use a printer named "Printer_POS_80". You can change this in `main.go`:

```go
printerName := "Your_Printer_Name"
```

### Image Settings
- **Quality**: 90% (configurable in `htmltoimage.go`)
- **Format**: PNG (configurable in `htmltoimage.go`)
- **Full Page**: Captures entire HTML content
- **Timeout**: 30 seconds (configurable in `htmltoimage.go`)

### Output Settings
- **Output Directory**: `./tmp/printy/`
- **Chrome Data**: `./tmp/chrome-data/`
- **Temporary Files**: Automatically cleaned up

## Modules

### HTMLToImageConverter
- Converts HTML content directly to images using ChromeDP
- Reuses Chrome instance for faster subsequent conversions
- Supports full page capture
- Thread-safe implementation

### ImagePrinter
- Prints images using CUPS
- Supports printer options
- Includes printer discovery functionality
- Error handling for printer connectivity

## HTML Template Support

The system supports HTML templates with:
- **Tailwind CSS**: Via CDN (works offline with cached resources)
- **Custom CSS**: Regular CSS styling
- **Responsive Design**: Mobile-friendly layouts
- **Print Optimization**: Proper page sizing and margins

## Chrome Instance Management

The system includes intelligent Chrome instance management:
- **Instance Reuse**: Same Chrome process for multiple conversions
- **Thread Safety**: Safe for concurrent use
- **Clean Shutdown**: Proper resource cleanup
- **Error Recovery**: Handles Chrome crashes gracefully

## Development

### Adding New HTML Templates
1. Create HTML files in the `templates/` directory
2. Modify `loadHTMLFromTemplates()` in `htmltoimage.go` to load your template
3. The system will automatically use the new template

### Extending Functionality
- Add new image formats in `htmltoimage.go`
- Implement additional printer options in `imageprinter.go`
- Add support for different Chrome configurations
- Integrate with template engines (Go templates, etc.)

## Troubleshooting

### Common Issues

1. **"Chrome not found"**
   - Install Chrome/Chromium using the instructions above
   - Ensure it's in your PATH

2. **"Failed to print image"**
   - Check that your printer is configured in CUPS
   - Verify the printer name in the configuration
   - Ensure the printer is online and accessible

3. **"ChromeDP screenshot failed"**
   - Check if Chrome is running properly
   - Verify the HTML template is valid
   - Check Chrome data directory permissions

### Debug Mode
To debug conversion issues, you can modify the temporary file cleanup to keep files for inspection:

```go
// Comment out the defer cleanup to keep files
// defer func() {
//     os.Remove(imagePath)
// }()
```

### Chrome Instance Debugging
To debug Chrome instance issues:
1. Check the `./tmp/chrome-data/` directory
2. Look for Chrome process logs
3. Verify Chrome can start manually

## Performance

- **First Conversion**: ~2-3 seconds (Chrome startup)
- **Subsequent Conversions**: ~0.5-1 second (reused Chrome instance)
- **Memory Usage**: ~100-200MB (Chrome process)
- **CPU Usage**: Low when idle, high during conversion

## License

This project is a proof of concept for educational purposes.