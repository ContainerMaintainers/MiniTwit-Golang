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

	who := uint(user) // SESSION USER ID

	whom, err := getUserId(username)
	if err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(404)
		return
	}

	follow := entities.Follower{
		Who_ID:  who, 
		Whom_ID: whom,
	}

	err = database.DB.Create(&follow).Error
	if err != nil { //when user is already following 'whom'
		log.Print("Bad request during " + c.Request.RequestURI + ": " + "Already following " + username)
		c.Status(400)
		return
	}

	c.String(200, fmt.Sprintf("You are now following \"%s\"", username))

}

// ENDPOINT: GET /:username
func username(c *gin.Context) { // Displays a user's tweets

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
			//followed := GetFollower(GetUser(username).ID, GetUser(user).ID)
			var users_page = false
			// If logged in is same as username endpoint
			if user == int(userID) {
				users_page = true
			}
			c.HTML(http.StatusOK, "timeline.html", gin.H{
				"title":         username + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"user":          username,
				//"followed":      followed,
				"user_page":     users_page,
				"messages":      GetMessages("myTimeline", user),
			})
		} else {
			// If not logged in
			c.HTML(http.StatusOK, "timeline.html", gin.H{
				"title":         username + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages("individual", int(userID)),
			})
		}
	} else {
		// if logged in
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":     "My Timeline",
			"user":      user,
			"private":   true,
			"user_page": true,
			"messages":  GetMessages("myTimeline", user),
		})
	}


	/*c.HTML(http.StatusOK, "timeline.html", gin.H{
		"messages": GetMessages("individual", int(userID)),
		"user": userID,
	})*/

}

// ENDPOINT: DELETE /:username/unfollow
func usernameUnfollow(c *gin.Context) { // Adds the current user as follower of the given user

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
		Who_ID:  who, 
		Whom_ID: whom,
	}

	err = database.DB.Where("Who_ID = ? AND Whom_ID = ?", unfollow.Who_ID, unfollow.Whom_ID).Delete(&unfollow).Error
	if err != nil { // when user is already following 'whom'
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.Status(400)
		return
	}

	c.String(200, fmt.Sprintf("You are no longer following \"%s\"", username)) // Had to make it 200 to satisfy tests for some reason

}

// ENDPOINT: POST /login
func login_user(c *gin.Context) { //Logs the user in.

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

	// if POST req?
	if _, err := getUserId(body.Username); err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": Invalid username " + body.Username)
		error = "Invalid username"
	} else if _, err := checkPasswordHash(body.Username, body.Password); err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error() + ", invalid password")
		error = "Invalid password"
	}

	if error == "" {
		// give message "You were logged in."
		// set session user to body.Username

		if userID, err := getUserId(body.Username); err != nil {
			log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
			user = -1
		} else {
			user = int(userID)
		}

		// redirect to timeline ("/")
		//user_path := "/" + body.Username
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
		//c.SetCookie("user", body.Username, 3600, "/", "/", false, false)
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

// Registration Helper Function
func ValidRegistration(c *gin.Context, username string, email string, password1 string, password2 string) bool {
	
	//error = ""
	if password1 == "" || email == "" || username == "" {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Missing Field")
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "All fields are required",
		})
		return false
	} else if email == "" || !strings.Contains(email, "@") {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Invalid email")
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "You have to enter a valid email address",
		})
		return false
	} else if password1 != password2 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Passwords do not match")
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "The two passwords do not match",
		})
		return false
	} else if _, err := getUserId(username); err == nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + " Username already taken")
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "The username is already taken",
		})
		return false
	}
	return true
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

	if ValidRegistration(c, body.Username, body.Email, body.Password, body.Password2) {
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

func getUserId(username string) (uint, error) { // Convenience method to look up the id for a username.

	var user entities.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}
