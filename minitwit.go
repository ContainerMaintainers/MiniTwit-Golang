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

func public(c *gin.Context) { //Displays the latest messages of all users

	c.JSON(200, gin.H{
		"message": "pong",
	})

	/* from python project:
	   return render_template('timeline.html', messages=query_db('''
	       select message.*, user.* from message, user
	       where message.flagged = 0 and message.author_id = user.user_id
	       order by message.pub_date desc limit ?''', [PER_PAGE]))
	*/
}

func username(c *gin.Context) { //Displays a user's tweets

	username := c.Param("username") //gets the <username> from the url
	c.JSON(200, gin.H{
		"username": username,
	})

}

func usernameFollow(c *gin.Context) { //Adds the current user as follower of the given user
	username := c.Param("username")
	c.JSON(200, gin.H{
		"username": username,
	})
}

// SIM ENDPOINTS:

func simLatest(c *gin.Context) {
	c.String(200, "simLatest")
}

func simRegister(c *gin.Context) {
	c.String(200, "simRegister")
}

func simMsgs(c *gin.Context) {
	c.String(200, "simMsgs")
}

func simPostUserMsg(c *gin.Context) {
	c.String(200, "simPostUserMsg")
}

func simGetUserMsg(c *gin.Context) {
	c.String(200, "simGetUserMsg")
}

func simGetUserFllws(c *gin.Context) {
	c.String(200, "simGetUserFllws")
}

func simPostUserFllws(c *gin.Context) {
	c.String(200, "simPostUserFllws")
}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", ping)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", usernameFollow)

	// SIM ENDPOINTS:

	router.GET("/sim/latest", simLatest)
	router.POST("/sim/register", simRegister)
	router.GET("/sim/msgs", simMsgs)
	router.POST("/sim/msgs/:username", simPostUserMsg)
	router.GET("/sim/msgs/:username", simGetUserMsg)
	router.POST("/sim/fllws/:username", simPostUserFllws)
	router.GET("/sim/fllws/:username", simGetUserFllws)

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
