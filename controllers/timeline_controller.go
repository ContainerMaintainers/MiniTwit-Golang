package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ENDPOINT: GET /public
func Timeline(c *gin.Context) {
	page := c.DefaultQuery("page", "0")

	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, page),
		})
	} else {
		// if there exists a session user, show my timeline
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, page),
			"user":     user,
			"username": user_name,
		})
	}
}

// ENDPOINT: GET /:username
func UserTimeline(c *gin.Context) { // Displays an individual's timeline

	username := c.Param("username") // gets the <username> from the url
	userID, err := getUserId(username)
	page := c.DefaultQuery("page", "0")

	if err != nil {
		log.Print("Bad request during " + c.Request.RequestURI + ": " + " User " + username + " not found")
		c.Status(404)
		return
	}
	// if endpoint is a username
	if username != "" {
		// if logged in
		if user != -1 {
			followed := GetFollower(uint(userID), uint(user))
			var users_page = false
			// If logged in user == endpoint
			if user == int(userID) {
				users_page = true

				c.HTML(http.StatusOK, "timeline.html", gin.H{
					"title":     "My Timeline",
					"user":      user,
					"private":   true,
					"user_page": true,
					"messages":  GetMessages("myTimeline", user, page),
					"username":  username,
				})
			} else {
				// If following
				if followed == true {
					// If logged in and user != endpoint
					c.HTML(http.StatusOK, "timeline.html", gin.H{
						"title":         username + "'s Timeline",
						"user_timeline": true,
						"private":       true,
						"user":          username,
						"followed":      followed,
						"user_page":     users_page,
						"messages":      GetMessages("individual", int(userID), page),
					})
				} else {
					// If not following
					// If logged in and user != endpoint
					c.HTML(http.StatusOK, "timeline.html", gin.H{
						"title":         username + "'s Timeline",
						"user_timeline": true,
						"private":       true,
						"user":          username,
						"user_page":     users_page,
						"messages":      GetMessages("individual", int(userID), page),
					})
				}
			}
		} else {
			// If not logged in
			c.HTML(http.StatusOK, "timeline.html", gin.H{
				"title":         username + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages("individual", int(userID), page),
			})
		}
	}
}
