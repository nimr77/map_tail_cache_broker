package map_core

type MapResponse struct {
	MapRequest
	ProviderUrl string `json:"provider_url"`
	SavedUrl    string `json:"saved_url"`
}
