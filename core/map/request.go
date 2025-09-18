package map_core

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
