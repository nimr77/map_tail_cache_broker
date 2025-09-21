package storage_services

import (
	"context"
	"errors"
	"io"
	"log"
	"net/url"
	"strings"
)

// DownloadFile downloads a file from firebase storage
//
// url: the url of the file to download
//
// returns the file as a byte array
func DownloadFile(url string) ([]byte, error) {

	path, err := FromUrlToFileStoragePath(url)

	if err != nil {
		return nil, err
	}

	log.Printf("path %s", path)

	bucket, err := GetMapUploadingBucket()

	if err != nil {
		return nil, err
	}

	reader, err := bucket.Object(path).NewReader(context.Background())

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	bytes := make([]byte, reader.Attrs.Size)

	_, err = io.ReadFull(reader, bytes)

	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, errors.New("file is empty")
	}
	// if os.Getenv("env") == "debug" {
	// 	SaveFile(string(bytes), "downloaded_files", "jpg", "downloaded_file")
	// }

	return bytes, nil
}

// FromUrlToFileStoragePath converts a Firebase Storage URL (gs:// or download URL)
// to a file storage path.
func FromUrlToFileStoragePath(urlString string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		log.Printf("Error parsing URL %q: %v\n", urlString, err)
		return "", err
	}

	objectPath := parsedURL.Path // Get the path part of the URL

	// Clean up the object path for gs:// URLs and firebase download URLs
	if strings.HasPrefix(urlString, "gs://") {
		objectPath = strings.TrimPrefix(urlString, "gs://")
	} else if strings.HasPrefix(urlString, "https://firebasestorage.googleapis.com/v0/b/") {
		objectPath = strings.TrimPrefix(objectPath, "/o/") // Already have the path so just remove /o/
	}

	// URL decode the object path
	decodedPath, err := url.PathUnescape(objectPath)
	if err != nil {
		log.Printf("Error decoding object path %q: %v\n", objectPath, err)
		return "", err
	}

	return strings.ReplaceAll(decodedPath, "/download/storage/v1/b/ridehub-b57bb.firebasestorage.app/o/", ""), nil
}
