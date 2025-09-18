package map_core

import "fmt"

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

func (m MapRequest) GetMapTileFileName() string {
	return fmt.Sprintf("%s_%d_%d_%d.png", m.Provider, int(m.X), int(m.Y), int(m.Z))
}

func (m MapRequest) GetMapTilePath() string {
	return fmt.Sprintf("%s/%d/%f/%f/%f.png", m.Provider, m.ThemeMode, m.Z, m.X, m.Y)
}
