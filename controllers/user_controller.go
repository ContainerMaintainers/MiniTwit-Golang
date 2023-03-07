package controllers

import (
	"fmt"
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	latest = 0
	user   = -1
)

// ENDPOINT: POST /:username/follow
func usernameFollow(c *gin.Context) { //Adds the current user as follower of the given user

	if user == -1 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No user logged in")
		c.AbortWithStatus(401)
		return
	}

	username := c.Param("username")

	who := uint(user) // SHOULD GET SESSION USER ID

	whom, err := getUserId(username)
	if err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(404)
		return
	}

	follow := entities.Follower{
		Who_ID:  who, // !
		Whom_ID: whom,
	}

	err = database.DB.Create(&follow).Error
	if err != nil { //when user is already following 'whom'
		log.Print("Bad request during " + c.Request.RequestURI + ": " + "Already following " + username)
		c.Status(400)
		return
	}

	// c.JSON(200, gin.H{
	// 	"follower": follow,
	// })
	c.String(200, fmt.Sprintf("You are now following \"%s\"", username))

}

// ENDPOINT: GET /:username
func username(c *gin.Context) { //Displays a user's tweets

	username := c.Param("username") //gets the <username> from the url

	userID, err := getUserId(username)
	if err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " User " + username + " not found")
		c.Status(404)
		return
	}

	var messagesFromUser []entities.Message

	if err := database.DB.
		Where("author_id = ?", userID).
		Limit(Per_page).
		Find(&messagesFromUser).Error; err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.AbortWithStatus(404)
		return
	}

	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"messagesFromUser": messagesFromUser,
		"username":	username,
		"user": user,
	})

}

// ENDPOINT: DELETE /:username/unfollow
func usernameUnfollow(c *gin.Context) { //Adds the current user as follower of the given user

	if user == -1 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No user logged in")
		c.AbortWithStatus(401)
		return
	}

	username := c.Param("username")

	who := uint(user) // SHOULD GET SESSION USER ID

	whom, err := getUserId(username)
	if err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " User " + username + " not found")
		c.Status(404)
		return
	}

	unfollow := entities.Follower{
		Who_ID:  who, // !
		Whom_ID: whom,
	}

	err = database.DB.Where("Who_ID = ? AND Whom_ID = ?", unfollow.Who_ID, unfollow.Whom_ID).Delete(&unfollow).Error
	if err != nil { //when user is already following 'whom'
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(400)
		return
	}

	// c.JSON(204, gin.H{
	// 	"follower": unfollow,
	// })
	c.String(200, fmt.Sprintf("You are no longer following \"%s\"", username)) // Had to make it 200 to satisfy tests for some reason

}

// ENDPOINT: POST /login
func login_user(c *gin.Context) { //Logs the user in.
	//check if there exists a session user, if yes, redirect to timeline ("/")

	var body struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		log.Print("Ran into error when binding to context during " + c.Request.RequestURI + ": " + err.Error())
		c.AbortWithStatus(400)
		return
	}

	error := ""

	//if POST req?
	if _, err := getUserId(body.Username); err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": Invalid username " + body.Username)
		error = "Invalid username"
	} else if _, err := checkPasswordHash(body.Username, body.Password); err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error() + ", invalid password")
		error = "Invalid password"
	}

	if error == "" {
		//give message "You were logged in."
		//set session user to body.Username

		// Until session stuff is working, just keep track of the user through a global variable
		// In this case the id is replaced with the username
		if userID, err := getUserId(body.Username); err != nil {
			log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
			user = -1
		} else {
			user = int(userID)
		}

		//redirect to timeline ("/")
		//c.Redirect(200, "/")
		//user_path := "/" + body.Username
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
		c.SetCookie("user", body.Username, 3600, "/", "/", false, false)

		// Temporarily dont redirect
		//c.String(200, "You were logged in")

	} else {
		c.String(400, error)
	}

}

// ENDPOINT: GET /logout
func logout_user(c *gin.Context) {
	//clear session user
	user = -1
	
	//c.String(200, "You were logged out")
	c.Redirect(http.StatusFound, "/")
}

// ENDPOINT: GET /register
func register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"messages": "register page",
	})
}

// ENDPOINT: GET /login
func loginf(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"messages": "Login page",
	})
}

// ENDPOINT: POST /register
func register_user(c *gin.Context) {

	var body struct {
		Username  string `form:"username" json:"username"`
		Email     string `form:"email" json:"email"`
		Password  string `form:"password" json:"password"`
		Password2 string `form:"password2" json:"password2"`
	}

	err := c.Bind(&body)
	if err != nil {
		log.Print("Ran into error when binding to context during " + c.Request.RequestURI + ": " + err.Error())
		c.AbortWithStatus(400)
		return
	}

	error := ""

	if body.Username == "" {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No username provided")
		error = "You have to enter a username"
	} else if body.Email == "" || !strings.Contains(body.Email, "@") {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Invalid email")
		error = "You have to enter a valid email address"
	} else if body.Password == "" {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No password provided")
		error = "You have to enter a password"
	} else if body.Password != body.Password2 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Passwords do not match")
		error = "The two passwords do not match"
	} else if _, err := getUserId(body.Username); err == nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + " Username already taken")
		error = "The username is already taken"
	}

	if error == "" {
		user := entities.User{
			Username: body.Username,
			Email:    body.Email,
			Password: Salt_pwd(body.Password),
		}

		err := database.DB.Create(&user).Error
		if err != nil {
			log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
			c.AbortWithStatus(500)
			return
		}

		//c.String(200, "You were successfully registered and can login now")
		location := url.URL{Path: "/login"}
		c.Redirect(http.StatusFound, location.RequestURI())

	} else {
		c.String(400, error)
	}

}

func Salt_pwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Fatal("The given password could not be hashed. ", err)
	}

	return string(hash)
}

func checkPasswordHash(username string, enteredPW string) (bool, error) {
	var user entities.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(enteredPW)); err != nil {
		return false, err
	}

	if err := database.DB.Where("username = ? AND password = ?", username, user.Password).First(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}

func getUserId(username string) (uint, error) { //Convenience method to look up the id for a username.

	var user entities.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}
