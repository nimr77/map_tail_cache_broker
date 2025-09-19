package map_router

import (
	map_handler "map_broker/handlers/map"

	"github.com/gin-gonic/gin"
)

func GetRouter(r *gin.RouterGroup) {
	r.GET("/:x/:y/:z", map_handler.GetImageBaseOnXYZoom)
}
