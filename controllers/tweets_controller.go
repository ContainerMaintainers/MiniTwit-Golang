package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"

	"github.com/ContainerMaintainers/MiniTwit-Golang/monitoring"
	"github.com/gin-gonic/gin"
)

const Per_page int = 30

type JoinedMessage struct {
	Author_id     uint
	Username      string
	Text          string
	Pub_Date      uint
	FormattedDate string // new variable to hold the formatted date string
}

func LimitMessages(page string) (int, int) {
	messagesPerPage := 50
	p, err := strconv.Atoi(page)
	if err != nil {
		panic("Failed to parse page number")
	}
	offset := (p - 1) * messagesPerPage
	return offset, messagesPerPage
}

func GetMessages(timelineType string, user int, page string) []JoinedMessage {

	//offset:= 1
	//messagesPerPage := 1

	offset, messagesPerPage := LimitMessages(page)
	var joinedMessages []JoinedMessage

	if timelineType == "public" {
		// Join messages and users tables for public timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join users on users.id = messages.Author_id").
			Offset(offset).Limit(messagesPerPage).
			Order("messages.Pub_Date desc").
			Scan(&joinedMessages)

	} else if timelineType == "myTimeline" {
		// Join messages, users, and followers for my timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join followers on messages.author_id = followers.whom_id").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ? OR followers.who_id = ?)", false, user, user).
			Offset(offset).Limit(messagesPerPage).
			Order("messages.Pub_Date desc").
			Scan(&joinedMessages)

	} else if timelineType == "individual" {
		// Join messages and users for an individual's timeline
		database.DB.Table("messages").
			Select("messages.Author_id", "users.Username", "messages.Text", "messages.Pub_Date").
			Joins("left join users on users.id = messages.Author_id").
			Where("messages.flagged = ? AND (messages.author_id = ?)", false, user).
			Offset(offset).Limit(messagesPerPage).
			Order("messages.Pub_Date desc").
			Scan(&joinedMessages)
	}

	for i := range joinedMessages {
		timestamp := time.Unix(int64(joinedMessages[i].Pub_Date), 0)
		joinedMessages[i].FormattedDate = timestamp.Format("Jan 02, 2006 3:04pm")
	}

	return joinedMessages
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

	username := c.Param("username") // gets the <username> from the url
	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"title":     "My Timeline",
		"user":      user,
		"private":   true,
		"user_page": true,
		"messages":  GetMessages("myTimeline", user, "0"),
		"username":  username,
	})
}

// ENDPOINT: GET /ping
func Ping(c *gin.Context) {

	monitoring.CountEndpoint("/ping", "GET")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
