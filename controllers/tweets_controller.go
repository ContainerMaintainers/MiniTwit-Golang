package controllers

import (
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const Per_page int = 30

type JoinedMessage struct {
	Author_id uint
	Username  string
	Text	  string
	Pub_Date  uint
	}

func GetMessages(timelineType string, user int) []JoinedMessage {
	
	var joinedMessages []JoinedMessage
	
	if timelineType == "public" {
		// Join messages and users tables for public timeline
		database.DB.Table("messages").
		Select("messages.Author_id", "users.Username" ,"messages.Text", "messages.Pub_Date").
		Joins("left join users on users.id = messages.Author_id").Scan(&joinedMessages)

	} else if timelineType == "myTimeline" {
		// Join messages, users, and followers for my timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username" ,"messages.Text", "messages.Pub_Date").
			Joins("left join followers on messages.author_id = followers.whom_id").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ? OR followers.who_id = ?)", false, user, user).
			Scan(&joinedMessages)
	} else if timelineType == "individual"{
		// Join messages and users for an individual's timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username" ,"messages.Text", "messages.Pub_Date").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ?)", false, user).
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
	
	// if there is NO session user, show public timeline
	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user),
		})
	} else {
		// if there exists a session user, show my timeline
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("myTimeline", user),
			"user": user,
		})
	}
}

// ENDPOINT: GET /public
func public(c *gin.Context) {

	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user),
		})
	} else {
		// if there exists a session user, show my timeline
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user),
			"user": user,
		})
	}
}

// ENDPOINT: POST /add_message
func addMessage(c *gin.Context) { //Registers a new message for the user.

	if user == -1 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No user logged in")
		c.AbortWithStatus(401)
		return
	}

	var body struct {
		Text string `form:"text" json:"text"`
	}

	c.Bind(&body)

	message := entities.Message{
		Author_id: uint(user), // GET SESSION USER ID
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
	c.Redirect(http.StatusFound, "/")

}

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")

	router.GET("/ping", ping)
	router.GET("/", timeline)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.GET("/:username/follow", usernameFollow)
	router.DELETE("/:username/unfollow", usernameUnfollow)
	router.POST("/register", register_user)
	router.GET("/register", register)
	router.POST("/add_message", addMessage)
	router.POST("/login", login_user)
	router.GET("/login", loginf)
	router.GET("/logout", logout_user)

	router.GET("/sim/latest", simLatest)
	router.POST("/sim/register", simRegister)
	router.GET("/sim/msgs", simMsgs)
	router.POST("/sim/msgs/:username", simPostUserMsg)
	router.GET("/sim/msgs/:username", simGetUserMsg)
	router.POST("/sim/fllws/:username", simPostUserFllws)
	router.GET("/sim/fllws/:username", simGetUserFllws)

	return router

}
