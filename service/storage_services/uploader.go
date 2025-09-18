package storage_services

import (
	"context"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

// UploadFile uploads a file to firebase storage
//
// file: the file to upload
// fileHeader: the file header
// path: the path to upload the file to
//
// returns the url of the uploaded file
func UploadFile(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
	bucket, err := GetMapUploadingBucket()

	if err != nil {
		return "", err
	}

	writer := bucket.Object(path).NewWriter(context.Background())

	// writer.ACL = []storage.ACLRule{
	// 	{
	// 		Entity: storage.AllUsers,
	// 		Role:   storage.RoleReader,
	// 	},

	// }
	writer.ContentType = fileHeader.Header.Get("Content-Type")
	bytes := make([]byte, fileHeader.Size)
	file.Read(bytes)
	_, err = writer.Write(bytes)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	log.Println("uploaded to firebase storage")
	atter, err := bucket.Object(path).Attrs(context.Background())
	if err != nil {
		log.Println("error getting attrs", err.Error())
		return "", err
	}
	url := atter.MediaLink
	return url, nil
}

// UploadFile uploads a file to firebase storage publicly
//
// file: the file to upload
// fileHeader: the file header
// path: the path to upload the file to
//
// returns the url of the uploaded file
func UploadFilePublic(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
	bucket, err := GetMapUploadingBucket()

	if err != nil {
		return "", err
	}

	writer := bucket.Object(path).NewWriter(context.Background())

	writer.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	writer.ContentType = fileHeader.Header.Get("Content-Type")
	bytes := make([]byte, fileHeader.Size)
	file.Read(bytes)
	_, err = writer.Write(bytes)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	log.Println("uploaded to firebase storage")
	atter, err := bucket.Object(path).Attrs(context.Background())
	if err != nil {
		log.Println("error getting attrs", err.Error())
		return "", err
	}
	url := atter.MediaLink
	return url, nil
}
