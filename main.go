package main

import (
	"flag"
	"map_broker/router"
)

func main() {
	router.InitRouters()

	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	router.Router.Run(":" + *port)
}
