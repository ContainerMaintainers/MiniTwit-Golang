package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const Per_page int = 30

var latest = 0

func getUserId(username string) (uint, error) { //Convenience method to look up the id for a username.

	var user entities.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// ENDPOINT: /ping
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// ENDPOINT: /public
func public(c *gin.Context) { //Displays the latest messages of all users

	var messages []entities.Message

	if err := database.DB.Where("Flagged = false").Order("Pub_Date desc").Limit(Per_page).Find(&messages).Error; err != nil {
		c.AbortWithStatus(400)
	}

	c.JSON(200, gin.H{
		"messages": messages,
	})
}

// ENDPOINT: /:username
func username(c *gin.Context) { //Displays a user's tweets

	username := c.Param("username") //gets the <username> from the url

	userID, err := getUserId(username)

	if err != nil {
		c.Status(404)
		return
	}

	var messagesFromUser []entities.Message

	if err := database.DB.Where("author_id = ?", userID).Limit(Per_page).Find(&messagesFromUser).Error; err != nil {
		c.AbortWithStatus(404)
	}

	c.JSON(200, gin.H{
		"messagesFromUser": messagesFromUser,
	})

}

// ENDPOINT: /:username/follow
func usernameFollow(c *gin.Context) { //Adds the current user as follower of the given user

	//check if there exists a session user, if not, return error 401, try c.AbortWithStatus(401)

	username := c.Param("username")

	who := uint(4) // SHOULD GET SESSION USER ID

	whom, err := getUserId(username)

	if err != nil {
		c.Status(404)
		return
	}
	follow := entities.Follower{
		Who_ID:  who, // !
		Whom_ID: whom,
	}
	result := database.DB.Create(&follow)
	if result.Error != nil { //when user is already following 'whom'
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"follower": follow,
	})

}

// ENDPOINT: /:username/unfollow
func usernameUnfollow(c *gin.Context) { //Adds the current user as follower of the given user

	//check if there exists a session user, if not, return error 401, try c.AbortWithStatus(401)

	username := c.Param("username")

	who := uint(4) // SHOULD GET SESSION USER ID

	whom, err := getUserId(username)

	if err != nil {
		c.Status(404)
		return
	}
	unfollow := entities.Follower{
		Who_ID:  who, // !
		Whom_ID: whom,
	}
	result := database.DB.Where("Who_ID = ? AND Whom_ID = ?", unfollow.Who_ID, unfollow.Whom_ID).Delete(&unfollow)
	if result.Error != nil { //when user is already following 'whom'
		c.Status(400)
		return
	}

	c.JSON(204, gin.H{
		"follower": unfollow,
	})

}

// ENDPOINT: /add_message
func addMessage(c *gin.Context) { //Registers a new message for the user.

	//check if there exists a session user, if not, return error 401, try c.AbortWithStatus(401)

	var body struct {
		Text string `json:"text"`
	}

	c.BindJSON(&body)

	message := entities.Message{
		Author_id: 4, // AUTHOR ID SHOULD GET SESSION USER ID
		Text:      body.Text,
		Pub_Date:  uint(time.Now().Unix()),
		Flagged:   false,
	}

	result := database.DB.Create(&message)

	if result.Error != nil {
		c.Status(400)
		return
	}

	//redirect to timeline ("/")
	c.Redirect(200, "/") // For some reason, this returns error 500, but I assume it's because the path doesn't exist yet?
}

// ENDPOINT: /register
func register(c *gin.Context) {

	var body struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Password2 string `json:"password2"`
		Email     string `json:"email"`
	}

	c.BindJSON(&body)

	error := ""

	if body.Username == "" {
		error = "You have to enter a username"
	} else if body.Email == "" || !strings.Contains(body.Email, "@") {
		error = "You have to enter a valid email address"
	} else if body.Password == "" {
		error = "You have to enter a password"
	} else if body.Password != body.Password2 {
		error = "The two passwords do not match"
	} else if id, _ := getUserId(body.Username); id != 0 {
		error = "The username is already taken"
	} else {
		user := entities.User{
			Username: body.Username,
			PW_Hash:  body.Password, // UPDATE SO PASSWORD IS HASHED
			Email:    body.Email,
		}

		database.DB.Create(&user)

	}

	c.String(200, error)

}

// SIM ENDPOINTS & HELPER FUNCTIONS:

func notReqFromSimulator(request *http.Request) gin.H {
	from_simulator := request.Header.Get("Authorization")
	if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		return gin.H{"status": 403, "error_msg": "You are not authorized to use this resource!"}
	} else {
		return nil
	}
}

func updateLatest(request *http.Request) {
	latest_value, err := strconv.Atoi(request.Header.Get("latest"))
	if latest_value != -1 && err == nil {
		latest = latest_value
	}
}

// ENDPOINT: GET /sim/latest
func simLatest(c *gin.Context) {
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
	} else if id, _ := getUserId(body.Username); id != 0 {
		error = "The username is already taken"
	} else {
		user := entities.User{
			Username: body.Username,
			PW_Hash:  body.Password, // UPDATE SO PASSWORD IS HASHED
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
		gorm.Model
		Text     string `json:"content"`
		Pub_Date uint   `json:"pub_date"`
		Username string `json:"user"`
	}

	var messages []MessageUser

	num_of_msgs, err := strconv.Atoi(c.Request.Header.Get("no"))
	if err != nil {
		log.Fatal(err)
	}

	if err := database.DB.Table("messages").
		Joins("join users on messages.author_id = users.id").
		Where("messages.flagged = ?", false).Order("messages.pub_date desc").
		Limit(num_of_msgs).Find(&messages); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
		gorm.Model
		Text     string `json:"content"`
		Pub_Date uint   `json:"pub_date"`
		Username string `json:"user"`
	}

	var messages []MessageUser

	num_of_msgs, err := strconv.Atoi(c.Request.Header.Get("no"))
	if err != nil {
		log.Fatal(err)
	}

	if err := database.DB.Table("messages").
		Joins("join users on messages.author_id = users.id").
		Where("messages.flagged = ? AND users.username = ?", false, username).Order("messages.pub_date desc").
		Limit(num_of_msgs).Find(&messages); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		c.AbortWithStatus(404)
	}

	num_of_followers, err := strconv.Atoi(c.Request.Header.Get("no"))
	if err != nil {
		num_of_followers = 100
	}

	type Username struct {
		Username string
	}

	var usernames []Username

	if err := database.DB.Table("followers").
		Joins("join users on followers.whom_id = users.id").
		Where("followers.who_id = ?", user_id).
		Limit(num_of_followers).Find(&usernames); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		c.AbortWithStatus(404)
	}

	if body.Follow != "" {

		follow_user_id, err := getUserId(body.Follow)
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(404)
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
			log.Fatal(err)
			c.AbortWithStatus(404)
		}

		follower := entities.Follower{
			Who_ID:  user_id,
			Whom_ID: unfollow_user_id,
		}

		database.DB.Delete(&follower)

		c.String(204, "")
	}

	c.String(400, "")

}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", ping)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", usernameFollow)
	router.DELETE("/:username/unfollow", usernameUnfollow)
	router.POST("/register", register)
	router.POST("/add_message", addMessage)

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
}

func main() {
	database.ConnectToTestDatabase()
	database.MigrateEntities()

	router := setupRouter()
	router.Run() // port 8080
}
