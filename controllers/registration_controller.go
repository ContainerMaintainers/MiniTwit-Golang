package controllers

import (
	//"fmt"
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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
	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(enteredPW)); err != nil {
		return false, err
	} else if err := database.DB.Where("username = ? AND password = ?", username, user.Password).First(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}

// ENDPOINT: GET /register
func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"messages": "register page",
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

		location := url.URL{Path: "/login"}
		c.Redirect(http.StatusFound, location.RequestURI())

	} else {
		c.String(400, error)
	}
}