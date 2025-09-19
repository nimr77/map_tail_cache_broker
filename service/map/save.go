package map_services

import (
	map_core "map_broker/core/map"
	"map_broker/service/storage_services"
)

func UploadMapTileToStorageService(request map_core.MapRequest, fileBytes []byte) (string, error) {
	path := request.GetMapTilePath()

	url, err := storage_services.UploadFileBytes(fileBytes, "image/jpeg", path)

	if err != nil {
		return "", err
	}

	return url, nil
}
