package controllers

import (
	//"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// login helper function
func valid_login(c *gin.Context, message string, user int, username string) bool {
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
		// if valid login, direct to user's timeline "/username"
		user_name = username
		c.Redirect(http.StatusFound, "/"+username)
	}
	return true

}

// ENDPOINT: POST /login
func Login_user(c *gin.Context) { //Logs the user in.

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
		valid_login(c, "Invalid username", -1, body.Username)
	} else if _, err := checkPasswordHash(body.Username, body.Password); err != nil {
		// if invalid password
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error() + ", invalid password")
		valid_login(c, "Invalid password", -1, body.Username)
	} else {
		if userID, err := getUserId(body.Username); err != nil {
			log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
			user = -1
		} else {
			// set session user
			user = int(userID)
			valid_login(c, "You are logged in!", user, body.Username)
		}
	}
}

// ENDPOINT: GET /login
func Loginf(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"messages": "Login page",
	})
}

// ENDPOINT: GET /logout
func logout_user(c *gin.Context) {
	//clear session user
	user = -1

	c.Redirect(http.StatusFound, "/")
}
