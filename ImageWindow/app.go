package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Windows API imports
var (
	user32          = syscall.NewLazyDLL("user32.dll")
	getDpiForSystem = user32.NewProc("GetDpiForSystem")
)

// App struct
type App struct {
	ctx          context.Context
	originalSize struct {
		width  int
		height int
	}
	dpiScale float64
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		dpiScale: 1.0,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Store initial window size
	width, height := runtime.WindowGetSize(ctx)
	a.originalSize.width = width
	a.originalSize.height = height

	// Initialize DPI scaling
	a.initDpiScale()
}

// Initialize DPI scaling factor
func (a *App) initDpiScale() {
	const defaultDPI = 96.0

	// Get system DPI
	var dpi uint32 = defaultDPI
	ret, _, _ := getDpiForSystem.Call()
	if ret != 0 {
		dpi = uint32(ret)
	}

	// Calculate scaling factor
	a.dpiScale = float64(dpi) / defaultDPI
	fmt.Printf("System DPI: %d, Scaling factor: %.2f\n", dpi, a.dpiScale)
}

// OpenImageDialog opens a file selection dialog
func (a *App) OpenImageDialog() string {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Image",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Image Files (*.png;*.jpg;*.jpeg;*.gif;*.bmp)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.gif;*.bmp",
			},
		},
	})

	if err != nil || filePath == "" {
		return ""
	}

	return filePath
}

// GetImageDimensions gets the image dimensions
func (a *App) GetImageDimensions(filePath string) map[string]int {
	file, err := os.Open(filePath)
	if err != nil {
		return map[string]int{"width": 0, "height": 0}
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return map[string]int{"width": 0, "height": 0}
	}

	return map[string]int{"width": img.Width, "height": img.Height}
}

// GetImageBase64 converts image to base64 string
func (a *App) GetImageBase64(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	// Determine MIME type based on file extension
	mimeType := "image/jpeg" // default
	ext := filepath.Ext(filePath)
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".bmp":
		mimeType = "image/bmp"
	}

	// Create data URL
	base64Data := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)
}

// ResizeWindow adjusts window size to match image dimensions
func (a *App) ResizeWindow(width, height int) {
	// Adjust window size according to DPI scaling, using ceiling
	scaledWidth := int(float64(width)/a.dpiScale + 0.999)   // ceiling
	scaledHeight := int(float64(height)/a.dpiScale + 0.999) // ceiling

	fmt.Printf("Image dimensions: %dx%d, DPI scaling: %.2f, Adjusted window size: %dx%d\n",
		width, height, a.dpiScale, scaledWidth, scaledHeight)

	// Apply window size
	runtime.WindowSetSize(a.ctx, scaledWidth, scaledHeight)
	runtime.WindowCenter(a.ctx)
}

// RestoreWindowSize restores original window size
func (a *App) RestoreWindowSize() {
	runtime.WindowSetSize(a.ctx, a.originalSize.width, a.originalSize.height)
	runtime.WindowCenter(a.ctx)
}

// ExitApp quits the application
func (a *App) ExitApp() {
	runtime.Quit(a.ctx)
}
