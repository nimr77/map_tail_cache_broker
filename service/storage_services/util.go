package storage_services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

func GetClient() (*storage.Client, error) {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetMapUploadingBucket() (*storage.BucketHandle, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.Bucket(" map-cached"), nil
}

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
