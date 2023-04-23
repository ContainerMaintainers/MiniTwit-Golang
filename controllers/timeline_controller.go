package controllers

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// ENDPOINT: GET /
func timeline(c *gin.Context) {

	// if there is NO session user, show public timeline
	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, 0),
		})
	} else {
		// if there exists a session user, show my timeline
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("myTimeline", user, 0),
			"user":     user,
		})
	}
}

// ENDPOINT: GET /public
func public(c *gin.Context) {

	if user == -1 {
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, 0),
		})
	} else {
		// if there exists a session user, show my timeline
		username := c.Param("username")
		c.HTML(http.StatusOK, "timeline.html", gin.H{
			"messages": GetMessages("public", user, 0),
			"user":     user,
			"username": username,
		})
	}
}

// ENDPOINT: GET /:username
func username(c *gin.Context) { // Displays an individual's timeline

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
					"messages":  GetMessages("myTimeline", user, 0),
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
						"messages":      GetMessages("individual", int(userID), 0),
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
						"messages":      GetMessages("individual", int(userID), 0),
					})
				}
			}
		} else {
			// If not logged in
			c.HTML(http.StatusOK, "timeline.html", gin.H{
				"title":         username + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages("individual", int(userID), 0),
			})
		}
	}
}