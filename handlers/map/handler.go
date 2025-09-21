package map_handler

import (
	"log"
	map_core "map_broker/core/map"
	map_services "map_broker/service/map"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetImageBaseOnXYZoomHandler(c *gin.Context) {
	x := c.Param("x")
	y := c.Param("y")
	z := c.Param("z")

	log.Println("X:", x, "Y:", y, "Z:", z)

	provider := c.Query("provider")
	theme := c.Query("theme")

	log.Println("Provider:", provider, "Theme:", theme)

	themeInt := 0
	if theme == "light" {
		themeInt = 1
	} else {
		themeInt = 0
	}

	mapRequest := map_core.MapRequest{}

	// xFloat, err := strconv.ParseFloat(x, 64)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid x parameter"})
	// 	return
	// }

	// yFloat, err := strconv.ParseFloat(y, 64)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid y parameter"})
	// 	return
	// }

	// zFloat, err := strconv.ParseFloat(z, 64)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid z parameter"})
	// 	return
	// }

	mapRequest.X = x
	mapRequest.Y = y
	mapRequest.Z = z
	mapRequest.Provider = strings.ToLower(provider)
	mapRequest.ThemeMode = map_core.ThemeMode(themeInt)

	log.Printf("Map Request: %+v\n", mapRequest)

	imageBy, err := map_services.GetTail(mapRequest)
	if err != nil {
		log.Printf("Error getting image: %v", err.Error())
		c.JSON(400, gin.H{"error": "Failed to get image", "massage": err.Error()})
		return
	}

	c.Data(200, "image/png", imageBy)

}
