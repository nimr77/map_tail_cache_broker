package map_services

import map_core "map_broker/core/map"

func CheckIfMapTailExistInService(request map_core.MapRequest) (bool, string, error) {
	// check if the map tile exists in the service (e.g., local storage or cache)
	return false, "", nil
}
