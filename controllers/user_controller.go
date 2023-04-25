package controllers

import (
	//"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"

	"github.com/ContainerMaintainers/MiniTwit-Golang/monitoring"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	latest = 0
	user   = -1
)

// ENDPOINT: GET /:username/follow
func UsernameFollow(c *gin.Context) { // Adds the current user as follower of the given user

	monitoring.CountEndpoint("/:username/follow", "GET")

	if user == -1 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No user logged in")
		c.AbortWithStatus(401)
		return

	} else {

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
		//database.DB.Create(&follow)
		//c.String(200, fmt.Sprintf("You are now following \"%s\"", username))
		c.Redirect(http.StatusFound, "/"+username)
	}
}

func GetFollower(follower uint, following uint) bool {
	var follows []entities.Follower
	if follower == following {
		return false
	} else {
		database.DB.Find(&follows).Where("who_ID = ?", following).Where("whom_ID = ?", follower).First(&follows)
		return len(follows) > 0
	}
}

// ENDPOINT: GET /:username/unfollow
func UsernameUnfollow(c *gin.Context) { // Adds the current user as follower of the given user

	monitoring.CountEndpoint("/:username/unfollow", "GET")

	if user == -1 {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " No user logged in")
		c.AbortWithStatus(401)
		return
	} else {
		username := c.Param("username")
		who := uint(user) // SESSION USER ID
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

		//c.String(200, fmt.Sprintf("You are no longer following \"%s\"", username)) // Had to make it 200 to satisfy tests
		c.Redirect(http.StatusFound, "/"+username)
	}
}

// login helper function
func valid_login(c *gin.Context, message string, user int) bool {
	if message == "Invalid username" {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Invalid Username")
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Invalid Username",
		})
		return false
	} else if message == "Invalid password" {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " Invalid password")
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Invalid password",
		})
		return false
	} else if message == "You are logged in!" {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"title":     "My Timeline",
			"flashes":	 message,
			"user":      user,
			"private":   true,
			"user_page": true,
			"messages":  GetMessages("myTimeline", user, 0),
		})
		return true
	}
	return true

}

// ENDPOINT: POST /login
func Login_user(c *gin.Context) { //Logs the user in.

	monitoring.CountEndpoint("/login", "POST")

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

	if _, err := getUserId(body.Username); err != nil {
		// If invalid username
		log.Print("Bad request during " + c.Request.RequestURI + ": Invalid username " + body.Username)
		valid_login(c, "Invalid username", -1)
	} else if _, err := checkPasswordHash(body.Username, body.Password); err != nil {
		// if invalid password
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error() + ", invalid password")
		valid_login(c, "Invalid password", -1)
	} else {
		if userID, err := getUserId(body.Username); err != nil {
			log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
			user = -1
		} else {
			// set session user
			user = int(userID)
			valid_login(c, "You are logged in!", user)
		}
	}
}

// ENDPOINT: GET /logout
func Logout_user(c *gin.Context) {

	monitoring.CountEndpoint("/logout", "GET")

	//clear session user
	user = -1

	//c.String(200, "You were logged out")
	c.Redirect(http.StatusFound, "/")
}

// ENDPOINT: GET /register
func Register(c *gin.Context) {

	monitoring.CountEndpoint("/register", "GET")

	c.HTML(http.StatusOK, "register.html", gin.H{
		"messages": "register page",
	})
}

// ENDPOINT: GET /login
func Loginf(c *gin.Context) {

	monitoring.CountEndpoint("/login", "GET")

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
func Register_user(c *gin.Context) {

	monitoring.CountEndpoint("/register", "POST")

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

