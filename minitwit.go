package main

import (
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	database.ConnectToDatabase()
	database.MigrateEntities()
}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // port 8080
}
