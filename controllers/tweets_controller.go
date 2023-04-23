package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"

	"github.com/ContainerMaintainers/MiniTwit-Golang/monitoring"
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
			Order("messages.Pub_Date desc").
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
			Order("messages.Pub_Date desc").
			Scan(&joinedMessages)

	} else if timelineType == "individual" {
		// Join messages and users for an individual's timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ?)", false, user).
			Offset(pagenum * Per_page).
			Limit(Per_page).
			Order("messages.Pub_Date desc").
			Scan(&joinedMessages)

	}

	return joinedMessages

}

// ENDPOINT: GET /ping
func Ping(c *gin.Context) {

	monitoring.CountEndpoint("/ping", "GET")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// ENDPOINT: GET /
func Timeline(c *gin.Context) {

	monitoring.CountEndpoint("/", "GET")

	// if there is NO session user, show public timeline
	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, 0),
		})
	} else {
		// if there exists a session user, show my timeline
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("myTimeline", user, 0),
			"user":     user,
		})
	}
}

// ENDPOINT: GET /public
func Public(c *gin.Context) {

	monitoring.CountEndpoint("/public", "GET")

	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, 0),
		})
	} else {
		// if there exists a session user, show my timeline
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, 0),
			"user":     user,
		})
	}
}

// ENDPOINT: GET /:username
func Username(c *gin.Context) { // Displays an individual's timeline

	monitoring.CountEndpoint("/:username", "GET")

	username := c.Param("username") // gets the <username> from the url
	userID, err := getUserId(username)
	if err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " User " + username + " not found")
		c.Status(404)
		return
	}
	// if endpoint is a username
	if username != "" {
		// if logged in
		if user != -1 {
			followed := GetFollower(uint(userID), uint(user))
			var users_page = false
			// If logged in user == endpoint
			if user == int(userID) {
				users_page = true

				c.HTML(http.StatusOK, "timeline.html", gin.H{
					"title":     "My Timeline",
					"user":      user,
					"private":   true,
					"user_page": true,
					"messages":  GetMessages("myTimeline", user, 0),
				})
			} else {
				// If following
				if followed == true {
					// If logged in and user != endpoint
					c.HTML(http.StatusOK, "timeline.html", gin.H{
						"title":         username + "'s Timeline",
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
					c.HTML(http.StatusOK, "timeline.html", gin.H{
						"title":         username + "'s Timeline",
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
			c.HTML(http.StatusOK, "timeline.html", gin.H{
				"title":         username + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages("individual", int(userID), 0),
			})
		}
	}
}

// ENDPOINT: POST /add_message
func AddMessage(c *gin.Context) { // Registers a new message for the user.

	monitoring.CountEndpoint("/add_message", "POST")

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

	c.Redirect(http.StatusFound, "/")

}
