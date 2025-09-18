package map_core

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

var mapTailsProviders = map[MapTailsProvider]MapProvider{
	MaptilerDark: {
		Name:    "maptiler",
		BaseUrl: "https://api.maptiler.com/maps/",
		Url:     "https://api.maptiler.com/maps/streets-v2-dark/{z}/{x}/{y}@2x.png?key=YOUR_MAPTILER_KEY",
	},
	MaptilerLight: {
		Name:    "maptiler",
		BaseUrl: "https://api.maptiler.com/maps/",
		Url:     "https://api.maptiler.com/maps/streets-v2/{z}/{x}/{y}@2x.png?key=YOUR_MAPTILER_KEY",
	},
}

func (m MapTailsProvider) GetMapTailsProvider() MapProvider {
	return mapTailsProviders[m]
}

func GetListOfMapTailsProviders() []MapProvider {
	providers := make([]MapProvider, len(mapTailsProviders))
	for _, provider := range mapTailsProviders {
		providers = append(providers, provider)
	}
	return providers
}
