package map_handler

import (
	"log"
	map_core "map_broker/core/map"

	"github.com/gin-gonic/gin"
)

func GetImageBaseOnXYZoom(c *gin.Context) {
	x := c.Param("x")
	y := c.Param("y")
	z := c.Param("z")

	log.Println("X:", x, "Y:", y, "Z:", z)

	provider := c.Query("provider")
	theme := c.Query("theme")

	log.Println("Provider:", provider, "Theme:", theme)

	themeInt := 0
	if theme == "light" {
		themeInt = 0
	} else {
		themeInt = 1
	}

	mapRequest := map_core.MapRequest{}
}
