package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Define command line arguments
	srcPtr := flag.String("src", "tasks", "Source directory path")
	tgtPtr := flag.String("tgt", "templates", "Target directory path")

	// Parse command line arguments
	flag.Parse()

	// Get directory paths from command line arguments or use defaults
	tasksDir := *srcPtr
	templatesDir := *tgtPtr

	// Ensure templates directory exists
	if err := ensureDir(templatesDir); err != nil {
		fmt.Printf("Failed to create templates directory: %v\n", err)
		waitForKeypress()
		return
	}

	// Count deleted and copied files
	deletedCount := 0
	copiedCount := 0

	// Create regex pattern to match PNG files with 'tpl' followed by 13 digits
	tplPattern := regexp.MustCompile(`^tpl\d{13}\.png$`)

	// Recursively traverse the tasks directory
	err := filepath.Walk(tasksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process PNG files
		if strings.ToLower(filepath.Ext(path)) == ".png" {
			filename := filepath.Base(path)

			// Delete if it's a tpl+13-digit PNG file
			if tplPattern.MatchString(filename) {
				if err := os.Remove(path); err != nil {
					fmt.Printf("Failed to delete %s: %v\n", path, err)
				} else {
					fmt.Printf("Deleted: %s\n", path)
					deletedCount++
				}
			} else {
				// Check if file exists in templates directory
				templatePath := filepath.Join(templatesDir, filename)
				if _, err := os.Stat(templatePath); os.IsNotExist(err) {
					// File doesn't exist in templates directory, copy it
					if err := copyFile(path, templatePath); err != nil {
						fmt.Printf("Failed to copy %s to %s: %v\n", path, templatePath, err)
					} else {
						fmt.Printf("Copied: %s to %s\n", path, templatePath)
						copiedCount++
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error traversing directory: %v\n", err)
	}

	fmt.Printf("\nOperation completed: %d files deleted, %d files copied\n", deletedCount, copiedCount)
	waitForKeypress()
}

// Ensure directory exists, create it if it doesn't
func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// Copy file from source to destination
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// Wait for user to press any key
func waitForKeypress() {
	fmt.Println("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
