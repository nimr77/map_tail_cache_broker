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
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
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
	url = strings.ReplaceAll(url, "{z}", fmt.Sprintf("%f", m.Z))
	url = strings.ReplaceAll(url, "{x}", fmt.Sprintf("%f", m.X))
	url = strings.ReplaceAll(url, "{y}", fmt.Sprintf("%f", m.Y))
	return url, nil
}

// func (m MapRequest) GetMapTileFileName() string {
// 	return fmt.Sprintf("%s_%d_%d_%d.png", m.Provider, int(m.X), int(m.Y), int(m.Z))
// }

func (m MapRequest) GetMapTilePath() string {
	return fmt.Sprintf("%s/%d/%f/%f/%f.png", m.Provider, m.ThemeMode, m.Z, m.X, m.Y)
}
