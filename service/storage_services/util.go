package storage_services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// SaveFile saves a string to a file within a specified folder.
func SaveFile(data string, folder string, ext string, filename string) error {
	// 1. Create the folder if it doesn't exist.
	err := os.MkdirAll(folder, 0755) // 0755: read/write/execute for owner, read/execute for others
	if err != nil {
		return fmt.Errorf("error creating folder: %w", err) // Wrap the error
	}

	// 2. Construct the full file path.
	filePath := filepath.Join(folder, filename+"."+ext)

	// 3. Create the file.
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close() // Ensure the file is closed when the function returns

	// 4. Write the data to the file.
	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	// 5. Optionally, sync the file to disk (less common, but increases durability).
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("error syncing file to disk: %w", err)
	}

	return nil
}

// DownloadFileFromUrl downloads a file from a URL and returns its content as a byte array.
func DownloadFileFromUrl(url string) ([]byte, error) {
	// 1. Make a GET request to the URL.
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer response.Body.Close()

	// 2. Check if the request was successful (status code 2xx).
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status code: %d", response.StatusCode)
	}

	// 3. Read all the bytes from the response body.
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}
