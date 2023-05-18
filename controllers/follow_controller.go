package controllers

import (
	//"fmt"
	"log"
	"net/http"

	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/monitoring"
	"github.com/gin-gonic/gin"
)

var (
	latest    = 0
	user      = -1
	user_name = ""
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

		c.Redirect(http.StatusFound, "/"+username)
	}
}

func getUserId(username string) (uint, error) { // Convenience method to look up the id for a username.

	var user entities.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}
