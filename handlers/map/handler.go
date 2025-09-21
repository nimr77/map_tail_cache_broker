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

	// // Parse and validate coordinates
	// xFloat, err := strconv.ParseFloat(x, 64)
	// if err != nil {
	// 	log.Printf("Invalid x parameter: %s", x)
	// 	c.JSON(400, gin.H{"error": "Invalid x parameter"})
	// 	return
	// }

	// yFloat, err := strconv.ParseFloat(y, 64)
	// if err != nil {
	// 	log.Printf("Invalid y parameter: %s", y)
	// 	c.JSON(400, gin.H{"error": "Invalid y parameter"})
	// 	return
	// }

	// zFloat, err := strconv.ParseFloat(z, 64)
	// if err != nil {
	// 	log.Printf("Invalid z parameter: %s", z)
	// 	c.JSON(400, gin.H{"error": "Invalid z parameter"})
	// 	return
	// }

	// // Validate coordinate ranges
	// if zFloat < 0 || zFloat > 18 {
	// 	log.Printf("Invalid zoom level: %f (must be 0-18)", zFloat)
	// 	c.JSON(400, gin.H{"error": "Invalid zoom level (must be 0-18)"})
	// 	return
	// }

	// maxTile := math.Pow(2, zFloat)
	// if xFloat < 0 || xFloat >= maxTile {
	// 	log.Printf("Invalid x coordinate: %f (must be 0-%f)", xFloat, maxTile-1)
	// 	c.JSON(400, gin.H{"error": "Invalid x coordinate"})
	// 	return
	// }

	// if yFloat < 0 || yFloat >= maxTile {
	// 	log.Printf("Invalid y coordinate: %f (must be 0-%f)", yFloat, maxTile-1)
	// 	c.JSON(400, gin.H{"error": "Invalid y coordinate"})
	// 	return
	// }

	// // Convert XYZ coordinates to TMS coordinates for MapTiler
	// // XYZ has Y=0 at top, TMS has Y=0 at bottom
	// tmsY := maxTile - 1 - yFloat

	// log.Printf("Original XYZ coordinates: X=%f, Y=%f, Z=%f", xFloat, yFloat, zFloat)
	// log.Printf("Converted TMS coordinates: X=%f, Y=%f, Z=%f", xFloat, tmsY, zFloat)

	// mapRequest.X = strconv.FormatFloat(xFloat, 'f', -1, 64)
	// mapRequest.Y = strconv.FormatFloat(tmsY, 'f', -1, 64)
	// mapRequest.Z = strconv.FormatFloat(zFloat, 'f', -1, 64)

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
