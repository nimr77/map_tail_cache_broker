package map_core

import "os"

type MapTailsProvider int

const (
	MaptilerDark  MapTailsProvider = 0
	MaptilerLight MapTailsProvider = 1
)

type MapProvider struct {
	Name    string
	BaseUrl string
	Url     string
}

func getApiKey() string {
	return os.Getenv("MAPTILER_API_KEY")
}

func getMapTailsProviders() map[MapTailsProvider]MapProvider {
	apiKey := getApiKey()
	return map[MapTailsProvider]MapProvider{
		MaptilerDark: {
			Name:    "maptiler",
			BaseUrl: "https://api.maptiler.com/maps/",
			Url:     "https://api.maptiler.com/maps/streets-v2-dark/{z}/{x}/{y}@2x.png?key=" + apiKey,
		},
		MaptilerLight: {
			Name:    "maptiler",
			BaseUrl: "https://api.maptiler.com/maps/",
			Url:     "https://api.maptiler.com/maps/streets-v2/{z}/{x}/{y}@2x.png?key=" + apiKey,
		},
	}
}

func (m MapTailsProvider) GetMapTailsProvider() MapProvider {
	return getMapTailsProviders()[m]
}

func GetListOfMapTailsProviders() []MapProvider {
	mapTailsProviders := getMapTailsProviders()
	providers := make([]MapProvider, 0, len(mapTailsProviders))
	for _, provider := range mapTailsProviders {
		providers = append(providers, provider)
	}
	return providers
}
