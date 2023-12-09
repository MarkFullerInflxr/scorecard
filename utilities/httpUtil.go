package utils

import (
	"net/http"
	"os"
)

func GetMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the first 512 bytes to determine the file type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Get the MIME type based on the content
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
