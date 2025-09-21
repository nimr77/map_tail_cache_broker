package storage_services

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func GetClient() (*storage.Client, error) {
	// Check if we have a service account key file
	credentialsFile := "map-broker-jaywalk-75c83aba05cf.json"

	var client *storage.Client
	var err error

	if _, err := os.Stat(credentialsFile); err == nil {
		// Use the service account key file
		client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(credentialsFile))
	} else {
		// Fall back to default credentials (for development)
		client, err = storage.NewClient(context.Background())
	}

	return client, err
}

func GetMapUploadingBucket() (*storage.BucketHandle, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.Bucket("map-cached"), nil
}
