package router

import (
	map_router "map_broker/router/map"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func InitRouters() {

	map_router.GetRouter(Router.Group("/map"))
}
