package map_core

import (
	"errors"
	"fmt"
	"strings"
)

type ThemeMode int

const (
	ThemeModeDark  ThemeMode = 0
	ThemeModeLight ThemeMode = 1
)

type MapRequestMeta struct {
	Provider  string    `json:"provider"`
	ThemeMode ThemeMode `json:"themeMode"`
}

type MapRequest struct {
	MapRequestMeta
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

func (m MapRequest) GetMapProvider() (MapTailsProvider, error) {
	if m.Provider == "maptiler" {
		if m.ThemeMode == ThemeModeDark {
			return MaptilerDark, nil
		} else {
			return MaptilerLight, nil
		}
	}

	return MapTailsProvider(0), errors.New("unknown map provider")
}

func (m MapRequest) GetFullMapTailUrl() (string, error) {
	provider, err := m.GetMapProvider()

	if err != nil {
		return "", err
	}

	url := provider.GetMapTailsProvider().Url

	// Check if API key is missing
	if strings.Contains(url, "key=") && strings.HasSuffix(url, "key=") {
		return "", errors.New("MAPTILER_API_KEY environment variable is not set")
	}

	url = strings.ReplaceAll(url, "{z}", fmt.Sprintf("%s", m.Z))
	url = strings.ReplaceAll(url, "{x}", fmt.Sprintf("%s", m.X))
	url = strings.ReplaceAll(url, "{y}", fmt.Sprintf("%s", m.Y))
	return url, nil
}

// func (m MapRequest) GetMapTileFileName() string {
// 	return fmt.Sprintf("%s_%d_%d_%d.png", m.Provider, int(m.X), int(m.Y), int(m.Z))
// }

func (m MapRequest) GetMapTilePath() string {
	return fmt.Sprintf("%s/%d/%s/%s/%s.png", m.Provider, m.ThemeMode, m.Z, m.X, m.Y)
}
