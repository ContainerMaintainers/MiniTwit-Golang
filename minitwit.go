package main

import (
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", ping)
	return router

}

func init() {
	initializers.LoadEnvVars()
	database.ConnectToDatabase()
	database.MigrateEntities()
}

func main() {

	router := setupRouter()
	router.Run() // port 8080
}
