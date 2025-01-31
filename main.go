package main

import (
	"github.com/gin-gonic/gin"
	"github.com/samuelIkoli/GODM/routes"
)



func main(){

	server:= gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")

}