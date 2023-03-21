package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/middleware"
	"github.com/gin-gonic/gin"
)

const Per_page int = 30

// ENDPOINT: GET /ping
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// ENDPOINT: GET /public
func public(c *gin.Context) { //Displays the latest messages of all users

	var messages []entities.Message

	if err := database.DB.Where("Flagged = false").Order("Pub_Date desc").Limit(Per_page).Find(&messages).Error; err != nil {
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
		c.AbortWithStatus(400)
		return
	}

	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"messages": messages,
	})
}

// ENDPOINT: GET /
func timeline(c *gin.Context) {

	// check if there exists a session user, if not, return all messages
	// For now just reuse the same endpoint handler as /public

	user, _ := c.Get("user")
	// if !exists {
	// 	log.Println("i dont exist")
	// }
	// log.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	// log.Println(user)
	// if user == nil {
	// 	c.Redirect(301, "/public")
	// }

	// usr := user.(entities.User)
	// log.Println("cccccccccccccccccccccccccccccccccccccccccccccccc")

	// tokenString, err := c.Cookie("UserAuthorization")
	// claims := jwt.MapClaims{}
	// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("SECRET_FOR_JWT_TOKEN"), nil
	// })
	// if err != nil {
	// 	log.Println("aaaaaaaaaaaaa")
	// 	log.Println(claims["sub"])
	// 	c.AbortWithStatus(403)
	// 	return
	// }
	// log.Println("bbbbbbbbbbbbbb")
	// if !token.Valid {
	// 	c.Redirect(301, "/public")
	// }
	// if usr.ID == 0 { // when it fetched no user
	// 	log.Println("dddddddddddddddddddddddddddddddddddddddddddd")
	// 	c.Redirect(301, "/public")
	// }

	var messages []entities.Message

	if err := database.DB.Table("messages").
		Joins("left join followers on messages.author_id = followers.whom_id").
		Where("messages.flagged = ? AND (messages.author_id = ? OR followers.who_id = ?)", false, user, user).
		Limit(Per_page).Find(&messages).Error; err != nil { // ORDER BY DATE
		log.Print("Ran into error during " + c.Request.RequestURI + ": " + err.Error())
	}

	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"messages": messages,
		"user":     user,
	})
}

// ENDPOINT: POST /add_message
func addMessage(c *gin.Context) { //Registers a new message for the user.

	user, _ := c.Get("user")

	usr := user.(entities.User)

	// log.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	// log.Println(usr)
	// log.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	var body struct {
		Text string `form:"text" json:"text"`
	}

	c.Bind(&body)

	message := entities.Message{
		Author_id: uint(usr.ID), // AUTHOR ID SHOULD GET SESSION USER ID
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

	c.String(200, "Your message was recorded")

}

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")

	router.GET("/ping", ping)
	router.GET("/", middleware.RequireAuth, timeline)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", middleware.RequireAuth, usernameFollow)
	router.DELETE("/:username/unfollow", middleware.RequireAuth, usernameUnfollow)
	router.POST("/register", register_user)
	router.GET("/register", register)
	router.POST("/add_message", middleware.RequireAuth, addMessage)
	router.POST("/login", login_user)
	router.GET("/login", loginf)
	router.PUT("/logout", logoutf) // Changed temporarily to satisfy tests, should it be put or get?
	router.GET("/validate", middleware.RequireAuth, validate)

	router.GET("/sim/latest", simLatest)
	router.POST("/sim/register", simRegister)
	router.GET("/sim/msgs", simMsgs)
	router.POST("/sim/msgs/:username", simPostUserMsg)
	router.GET("/sim/msgs/:username", simGetUserMsg)
	router.POST("/sim/fllws/:username", simPostUserFllws)
	router.GET("/sim/fllws/:username", simGetUserFllws)

	return router

}
