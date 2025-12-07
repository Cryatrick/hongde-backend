package thirdparty

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
    "mime/multipart"

    "github.com/ryanbekhen/go-webp"
)

func SaveImageAsWebp(uploadedFile multipart.File,fileName string , uploadDir string) (string, error) {
	// Decode any image type
	img, _, err := image.Decode(uploadedFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Ensure upload directory exists
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%v.webp", fileName)
	fullPath := filepath.Join(uploadDir, filename)

	// Create destination file
	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

    // Encode to WebP
    if err := webp.Encode(img, 80.0, outFile); err != nil {
        return "", fmt.Errorf("webp encode failed: %w", err)
    }

	return filename, nil
}