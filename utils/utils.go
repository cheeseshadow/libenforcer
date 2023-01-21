package utils

import (
	"net/http"
	"os"
)

func GetFileContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)

	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func CheckIfFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else {
		return false
	}
}
