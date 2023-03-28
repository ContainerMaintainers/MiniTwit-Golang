package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/middleware"
	"github.com/gin-gonic/gin"
)

const Per_page int = 30

type JoinedMessage struct {
	Author_id uint
	Username  string
	Text      string
	Pub_Date  uint
}

func GetMessages(timelineType string, user int, pagenum int) []JoinedMessage {

	var joinedMessages []JoinedMessage

	if timelineType == "public" {
		// Join messages and users tables for public timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join users on users.id = messages.Author_id").
			Offset(pagenum * Per_page).
			Limit(Per_page).
			Scan(&joinedMessages)

	} else if timelineType == "myTimeline" {
		// Join messages, users, and followers for my timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join followers on messages.author_id = followers.whom_id").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ? OR followers.who_id = ?)", false, user, user).
			Offset(pagenum * Per_page).
			Limit(Per_page).
			Scan(&joinedMessages)

	} else if timelineType == "individual" {
		// Join messages and users for an individual's timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ?)", false, user).
			Offset(pagenum * Per_page).
			Limit(Per_page).
			Scan(&joinedMessages)

	}

	return joinedMessages

}

// ENDPOINT: GET /ping
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// ENDPOINT: GET /
func timeline(c *gin.Context) {

	user, _ := c.Get("user")

	usr := user.(entities.User)

	c.JSON(http.StatusOK, gin.H{
		"messages": GetMessages("myTimeline", int(usr.ID), 0),
		"user":     user,
	})
}

// ENDPOINT: GET /public
func public(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"messages": GetMessages("public", 0, 0),
	})
}

// ENDPOINT: GET /:username
func username(c *gin.Context) { // Displays an individual's timeline

	username := c.Param("username") // gets the <username> from the url
	userID, err := getUserId(username)
	if err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " User " + username + " not found")
		c.Status(404)
		return
	}

	user, _ := c.Get("user")

	usr := user.(entities.User)

	// if endpoint is a username
	if username != "" {
		// if logged in
		if usr.ID != 0 {
			followed := GetFollower(uint(userID), uint(usr.ID))
			var users_page = false
			// If logged in user == endpoint
			if user == int(userID) {
				users_page = true

				c.JSON(http.StatusOK, gin.H{
					"title":     "My Timeline ONE",
					"user":      user,
					"private":   true,
					"user_page": true,
					"messages":  GetMessages("myTimeline", int(usr.ID), 0),
				})
			} else {
				// If following
				if followed == true {
					// If logged in and user != endpoint
					c.JSON(http.StatusOK, gin.H{
						"title":         username + "'s Timeline TWO",
						"user_timeline": true,
						"private":       true,
						"user":          username,
						"followed":      followed,
						"user_page":     users_page,
						"messages":      GetMessages("individual", int(userID), 0),
					})
				} else {
					// If not following
					// If logged in and user != endpoint
					c.JSON(http.StatusOK, gin.H{
						"title":         username + "'s Timeline THREE",
						"user_timeline": true,
						"private":       true,
						"user":          username,
						"user_page":     users_page,
						"messages":      GetMessages("individual", int(userID), 0),
					})
				}
			}
		} else {
			// If not logged in
			c.JSON(http.StatusOK, gin.H{
				"title":         username + "'s Timeline FOUR",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages("individual", int(userID), 0),
			})
		}
	}
}

// ENDPOINT: POST /add_message
func addMessage(c *gin.Context) { //Registers a new message for the user.

	user, _ := c.Get("user")

	usr := user.(entities.User)

	var body struct {
		Text string `form:"text" json:"text"`
	}

	c.Bind(&body)

	message := entities.Message{
		Author_id: uint(usr.ID), // AUTHOR ID SHOULD GET SESSION USER ID
		Text:      body.Text,
		Pub_Date:  uint(time.Now().Unix()),
		Flagged:   false,
	}

	err := database.DB.Create(&message).Error
	if err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(400)
		return
	}

	//redirect to timeline ("/")
	//c.Redirect(200, "/") // For some reason, this returns error 500, but I assume it's because the path doesn't exist yet?
	// Temporarily dont redirect
	c.String(200, "Your message was recorded")

}

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")

	router.GET("/ping", ping)
	router.GET("/", middleware.RequireAuth, timeline)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", usernameFollow)
	router.DELETE("/:username/unfollow", usernameUnfollow)
	router.POST("/register", register_user)
	router.GET("/register", register)
	router.POST("/add_message", middleware.RequireAuth, addMessage)
	router.POST("/login", login_user)
	router.GET("/login", loginf)
	router.PUT("/logout", logoutf) // Changed temporarily to satisfy tests, should it be put or get?

	router.GET("/sim/latest", simLatest)
	router.POST("/sim/register", simRegister)
	router.GET("/sim/msgs", simMsgs)
	router.POST("/sim/msgs/:username", simPostUserMsg)
	router.GET("/sim/msgs/:username", simGetUserMsg)
	router.POST("/sim/fllws/:username", simPostUserFllws)
	router.GET("/sim/fllws/:username", simGetUserFllws)

	return router

}
