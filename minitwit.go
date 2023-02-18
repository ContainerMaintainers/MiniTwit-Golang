package main

import (
	"log"
	"strings"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"

	"github.com/gin-gonic/gin"
)

const Per_page int = 30

func getUserId(username string) (uint, error) { //Convenience method to look up the id for a username.

	var user entities.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func checkPasswordHash(username string, enteredPW string) (bool, error) {
	var user entities.User

	hashedEnteredPW := enteredPW //hash "enteredPW" with hash function we're using

	if err := database.DB.Where("username = ? AND pw_hash = ?", username, hashedEnteredPW).First(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func timeline(c *gin.Context) {
	c.JSON(200, gin.H{
		"everything": "yep",
	})
}

func public(c *gin.Context) { //Displays the latest messages of all users

	var messages []entities.Message

	if err := database.DB.Where("Flagged = false").Order("Pub_Date desc").Limit(Per_page).Find(&messages).Error; err != nil {
		c.AbortWithStatus(400)
	}

	c.JSON(200, gin.H{
		"messages": messages,
	})
}

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

func loginf(c *gin.Context) { //Logs the user in.
	//check if there exists a session user, if yes, redirect to timeline ("/")

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	c.BindJSON(&body)

	error := ""

	//if POST req?
	if _, err := getUserId(body.Username); err != nil {
		error = "Invalid username"
	} else if _, err := checkPasswordHash(body.Username, body.Password); err != nil {
		error = "Invalid password"
	} else {
		//give message "You were logged in."
		log.Printf("You were logged in")
		//set session user to body.Username

		//redirect to timeline ("/")
		c.Redirect(200, "/")

	}

	c.String(400, error)

}

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
	router.GET("/", timeline)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", usernameFollow)
	router.DELETE("/:username/unfollow", usernameUnfollow)
	router.POST("/register", register)
	router.POST("/add_message", addMessage)
	router.POST("/login", loginf)

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
	database.ConnectToDatabase()
	database.MigrateEntities()

	router := setupRouter()
	router.Run() // port 8080
}
