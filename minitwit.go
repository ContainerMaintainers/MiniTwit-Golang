package main

import (
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"

	"github.com/gin-gonic/gin"
)

// def get_user_id(username):
//     """Convenience method to look up the id for a username."""
//     rv = g.db.execute('select user_id from user where username = ?',
//                        [username]).fetchone()
//     return rv[0] if rv else None

func getUserId(username string) (uint, error) { //Convenience method to look up the id for a username.

	var result uint
	database.DB.Raw("select user_id from users where username = ?", username).Scan(&result)

	// if result == 0 {
	// 	return (error
	// }

	//return uint32(result)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func public(c *gin.Context) { //Displays the latest messages of all users

	c.JSON(200, gin.H{
		"message": "pong",
	})

	/* from python project:
	   return render_template('timeline.html', messages=query_db('''
	       select message.*, user.* from message, user
	       where message.flagged = 0 and message.author_id = user.user_id
	       order by message.pub_date desc limit ?''', [PER_PAGE]))
	*/
}

func username(c *gin.Context) { //Displays a user's tweets

	username := c.Param("username") //gets the <username> from the url
	c.JSON(200, gin.H{
		"username": username,
	})

}

func usernameFollow(c *gin.Context) { //Adds the current user as follower of the given user

	var body struct {
		Who_ID  uint
		Whom_ID uint
	}

	c.Bind(&body)

	//who  = set logged in user as who
	//whom = get id from username

	post := entities.Follower{
		Who_ID:  body.Who_ID,
		Whom_ID: body.Whom_ID,
	}

	result := database.DB.Create(&post)

	if result.Error != nil { //when user is already following 'whom' || or given wrong types in fields
		c.Status(400)
		return
	}

	//username := c.Param("username")
	c.JSON(200, gin.H{
		"follower": post,
	})

}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", ping)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", usernameFollow)
	return router

}

func init() {
	initializers.LoadEnvVars()
	database.ConnectToDatabase()
	database.MigrateEntities()
}

func main() {

	router := setupRouter()
	router.Run() // port 8080
}
