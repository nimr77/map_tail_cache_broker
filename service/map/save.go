package map_services

import (
	"log"
	map_core "map_broker/core/map"
	"map_broker/service/storage_services"
)

func UploadMapTileToStorageService(request map_core.MapRequest, fileBytes []byte) (string, error) {
	path := request.GetMapTilePath()

	url, err := storage_services.UploadFileBytes(fileBytes, "image/jpeg", path)

	if err != nil {
		log.Println("Error uploading map tile to storage service:", err.Error())
		return "", err
	}

	log.Println("Map tile uploaded to storage service at URL:", url)

	return url, nil
}
