package main

import (
	"strings"

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

	return uint(result), nil
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

func register(c *gin.Context) {

	var body struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Password2 string `json:"password2"`
		Email     string `json:"email"`
	}

	c.BindJSON(&body)

	error := ""

	if body.Username == "" {
		error = "You have to enter a username"
	} else if body.Email == "" || !strings.Contains(body.Email, "@") {
		error = "You have to enter a valid email address"
	} else if body.Password == "" {
		error = "You have to enter a password"
	} else if body.Password != body.Password2 {
		error = "The two passwords do not match"
	} else if id, _ := getUserId(body.Username); id != 0 {
		error = "The username is already taken"
	} else {
		user := entities.User{
			Username: body.Username,
			PW_Hash:  body.Password, // UPDATE SO PASSWORD IS HASHED
			Email:    body.Email,
		}

		database.DB.Create(&user)

	}

	c.String(200, error)

}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", ping)
	router.GET("/public", public)
	router.GET("/:username", username)
	router.POST("/:username/follow", usernameFollow)
	router.POST("/register", register)
	return router

}

func init() {
	initializers.LoadEnvVars()
}

func main() {
	database.ConnectToDatabase()
	database.MigrateEntities()

	router := setupRouter()
	router.Run() // port 8080
}
