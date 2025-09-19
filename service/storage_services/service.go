package storage_services

import (
	"context"

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
	return client.Bucket("map-cached"), nil
}
