package controllers

import (
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func notReqFromSimulator(request *http.Request) gin.H {
	from_simulator := request.Header.Get("Authorization")
	if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		return gin.H{"status": 403, "error_msg": "You are not authorized to use this resource!"}
	} else {
		return nil
	}
}

func updateLatest(request *http.Request) {
	log.Print("Updating latest")
	latest_value, err := strconv.Atoi(request.URL.Query().Get("latest"))
	if latest_value != -1 && err == nil {
		latest = latest_value
	} else if err != nil {
		log.Print("During updateLatest(): ", err)
		latest = -1
	}
}

// ENDPOINT: GET /sim/latest
func simLatest(c *gin.Context) {
	log.Print("/sim/latest ", latest)
	c.JSON(200, gin.H{"latest": latest})
}

// ENDPOINT: POST /sim/register
func simRegister(c *gin.Context) {

	updateLatest(c.Request)

	var body struct {
		Username string `json:"username"`
		Password string `json:"pwd"`
		Email    string `json:"email"`
	}

	c.BindJSON(&body)

	error := ""

	if body.Username == "" {
		error = "You have to enter a username"
	} else if body.Email == "" || !strings.Contains(body.Email, "@") {
		error = "You have to enter a valid email address"
	} else if body.Password == "" {
		error = "You have to enter a password"
	} else if _, err := getUserId(body.Username); err != nil {
		error = "The username is already taken"
	} else {
		user := entities.User{
			Username: body.Username,
			Password: Salt_pwd(body.Password),
			Email:    body.Email,
		}

		database.DB.Create(&user)

	}

	if error != "" {
		c.JSON(400, gin.H{"status": 400, "error_msg": error})
	} else {
		c.String(204, "")
	}
}

// ENDPOINT: GET /sim/msgs
func simMsgs(c *gin.Context) {

	updateLatest(c.Request)

	if notFromSimResponse := notReqFromSimulator(c.Request); notFromSimResponse != nil {
		c.JSON(403, notFromSimResponse)
		return
	}

	type MessageUser struct {
		ID       uint   `json:"ID"`
		Text     string `json:"content"`
		Pub_Date uint   `json:"pub_date"`
		Username string `json:"user"`
	}

	var messages []MessageUser

	num_of_msgs, err := strconv.Atoi(c.Request.URL.Query().Get("no"))
	if err != nil {
		log.Print("During /sim/msgs ", err)
		num_of_msgs = 100
	}

	if err := database.DB.Table("messages").
		Joins("join users on messages.author_id = users.id").
		Where("messages.flagged = ?", false).Order("messages.pub_date desc").
		Limit(num_of_msgs).
		Select("messages.ID, messages.text, messages.pub_date, users.username").
		Find(&messages).Error; err != nil {
		log.Print(err)
		c.AbortWithStatus(400)
		return
	}

	c.JSON(200, messages)
}

// ENDPOINT: POST /sim/msgs/:username
func simPostUserMsg(c *gin.Context) {

	var body struct {
		Content string `json:"content"`
	}

	c.BindJSON(&body)

	updateLatest(c.Request)

	if notFromSimResponse := notReqFromSimulator(c.Request); notFromSimResponse != nil {
		c.JSON(403, notFromSimResponse)
		return
	}

	username := c.Param("username")

	user_id, err := getUserId(username)
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(400)
		return
	}

	message := entities.Message{
		Author_id: user_id,
		Text:      body.Content,
		Pub_Date:  uint(time.Now().Unix()),
	}

	database.DB.Create(&message)

	c.String(204, "")

}

// ENDPOINT: GET /sim/msgs/:username
func simGetUserMsg(c *gin.Context) {

	updateLatest(c.Request)

	if notFromSimResponse := notReqFromSimulator(c.Request); notFromSimResponse != nil {
		c.JSON(403, notFromSimResponse)
		return
	}

	username := c.Param("username")

	type MessageUser struct {
		ID       uint   `json:"id"`
		Text     string `json:"content"`
		Pub_Date uint   `json:"pub_date"`
		Username string `json:"user"`
	}

	var messages []MessageUser

	num_of_msgs, err := strconv.Atoi(c.Request.URL.Query().Get("no"))
	if err != nil {
		log.Print("During /sim/msgs/:username ", err)
		c.AbortWithStatus(400)
		return
	}

	if err := database.DB.Table("messages").
		Joins("join users on messages.author_id = users.id").
		Where("messages.flagged = ? AND users.username = ?", false, username).Order("messages.pub_date desc").
		Limit(num_of_msgs).
		Select("messages.ID, messages.text, messages.pub_date, users.username").
		Find(&messages).Error; err != nil {
		log.Print(err)
		c.AbortWithStatus(400)
		return
	}

	c.JSON(200, messages)
}

// ENDPOINT: GET /sim/fllws/:username
func simGetUserFllws(c *gin.Context) {

	updateLatest(c.Request)

	if notFromSimResponse := notReqFromSimulator(c.Request); notFromSimResponse != nil {
		c.JSON(403, notFromSimResponse)
		return
	}

	username := c.Param("username")

	user_id, err := getUserId(username)
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(404)
		return
	}

	num_of_followers, err := strconv.Atoi(c.Request.URL.Query().Get("no"))
	if err != nil {
		num_of_followers = 100
		log.Print("During /sim/fllws/:username ", err)
	}

	type Username struct {
		Username string
	}

	var usernames []Username

	if err := database.DB.Table("users").
		Joins("join followers on followers.whom_id = users.id").
		Where("followers.who_id = ?", user_id).
		Limit(num_of_followers).Find(&usernames).Error; err != nil {
		log.Print(err)
	}

	var usernameStrings []string

	for _, username := range usernames {
		usernameStrings = append(usernameStrings, username.Username)
	}

	c.JSON(200, gin.H{"follows": usernameStrings})

}

// ENDPOINT: POST /sim/fllws/:username
func simPostUserFllws(c *gin.Context) {

	var body struct {
		Follow   string `json:"follow"`
		Unfollow string `json:"unfollow"`
	}

	c.BindJSON(&body)

	updateLatest(c.Request)

	if notFromSimResponse := notReqFromSimulator(c.Request); notFromSimResponse != nil {
		c.JSON(403, notFromSimResponse)
		return
	}

	username := c.Param("username")

	user_id, err := getUserId(username)
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(404)
		return
	}

	if body.Follow != "" {

		follow_user_id, err := getUserId(body.Follow)
		if err != nil {
			log.Print(err)
			c.AbortWithStatus(404)
			return
		}

		follower := entities.Follower{
			Who_ID:  user_id,
			Whom_ID: follow_user_id,
		}

		database.DB.Create(&follower)

		c.String(204, "")

	} else if body.Unfollow != "" {

		unfollow_user_id, err := getUserId(body.Unfollow)
		if err != nil {
			log.Print(err)
			c.AbortWithStatus(404)
			return
		}

		follower := entities.Follower{
			Who_ID:  user_id,
			Whom_ID: unfollow_user_id,
		}

		database.DB.Where("whom_id = ? and who_id = ?", follower.Whom_ID, follower.Who_ID).Delete(&follower)

		c.String(204, "")
	} else {
		c.String(400, "")
	}

}
