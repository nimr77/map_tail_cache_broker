package map_services

import (
	"context"
	"log"
	map_core "map_broker/core/map"
	"map_broker/service/storage_services"
)

func GetTail(request map_core.MapRequest) ([]byte, error) {
	exist, err := CheckIfMapTailExistInService(request)

	if err != nil {
		log.Println("Error checking if map tile exists in storage service:", err.Error())
		return nil, err
	}

	if exist {
		log.Println("Map tile exists in storage service, fetching from there.")
		return GetTailMapFromStorageService(request)
	}

	dowloadingUrl, err := request.GetFullMapTailUrl()

	if err != nil {
		log.Println("Error getting full map tile URL:", err.Error())
		return nil, err
	}

	by, err := storage_services.DownloadFileFromUrl(dowloadingUrl)

	if err != nil {
		log.Println("Error downloading map tile from URL:", err.Error())
		return nil, err
	}

	go UploadMapTileToStorageService(request, by)

	return by, nil

}

func CheckIfMapTailExistInService(request map_core.MapRequest) (bool, error) {
	// check if the map tile exists in the service (e.g., local storage or cache)

	path := request.GetMapTilePath()
	// if exists return true and the path
	// else return false and empty string
	// for now, we assume it does not exist
	// you can implement your own logic here to check in your storage service

	bucket, err := storage_services.GetMapUploadingBucket()

	if err != nil {
		return false, err
	}

	_, err = bucket.Object(path).Attrs(context.Background())

	if err == nil {
		return true, nil
	}

	return false, nil
}

func GetTailMapFromStorageService(request map_core.MapRequest) ([]byte, error) {
	path := request.GetMapTilePath()

	bytes, err := storage_services.DownloadFile(path)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
