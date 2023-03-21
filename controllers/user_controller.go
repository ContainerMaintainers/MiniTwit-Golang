package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	latest = 0
	//user   = -1
)

// ENDPOINT: POST /:username/follow
func usernameFollow(c *gin.Context) { //Adds the current user as follower of the given user

	user, _ := c.Get("user")

	who := user.(entities.User) // SHOULD GET SESSION USER ID

	username := c.Param("username")

	whom, err := getUserId(username)
	if err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(404)
		return
	}

	follow := entities.Follower{
		Who_ID:  who.ID,
		Whom_ID: whom,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + "Already following " + username)
		c.Status(400)
		return
	}

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

	if err := database.DB.Where("author_id = ?", userID).Limit(Per_page).Find(&messagesFromUser).Error; err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.AbortWithStatus(404)
		return
	}

	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"messagesFromUser": messagesFromUser,
	})

}

// ENDPOINT: DELETE /:username/unfollow
func usernameUnfollow(c *gin.Context) { //Adds the current user as follower of the given user

	user, _ := c.Get("user")

	who := user.(entities.User) // SHOULD GET SESSION USER ID

	username := c.Param("username")

	whom, err := getUserId(username)
	if err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " User " + username + " not found")
		c.Status(404)
		return
	}

	unfollow := entities.Follower{
		Who_ID:  who.ID,
		Whom_ID: whom,
	}

	err = database.DB.Where("Who_ID = ? AND Whom_ID = ?", unfollow.Who_ID, unfollow.Whom_ID).Delete(&unfollow).Error
	if err != nil { // what happens when trying to unfollow someone you're not currently following? idk, can't find minitwit py repo (:
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(400)
		return
	}

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
		//set session user to body.Username
		userId, _ := getUserId(body.Username) // we can '_' the error, since we check for that error earlier

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": userId,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_FOR_JWT_TOKEN")))

		if err != nil {
			c.JSON(401, gin.H{
				"error": "Failed creation of token",
			})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("UserAuthorization", tokenString, 3600*24*30, "", "", false, true) //set false to true when not on local host, 3600*24*30

		user_path := "/" + body.Username
		location := url.URL{Path: user_path}
		c.Redirect(http.StatusFound, location.RequestURI())

	} else {
		c.String(400, error)
	}

}

func validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(200, gin.H{
		"message": user,
	})
}

// ENDPOINT: PUT /logout
func logoutf(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode) //don't know what this line does
	c.SetCookie("UserAuthorization", "", 0, "", "", false, true)

	//c.Redirect(200, "/")
	// Temporarily don't redirect
	c.String(200, "You were logged out")
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
