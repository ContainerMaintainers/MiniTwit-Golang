package main

import (
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.loadEnvVars()
	initializers.ConnectToDatabase()
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
