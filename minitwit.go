package main

import (
	"flag"

	"github.com/ContainerMaintainers/MiniTwit-Golang/controllers"
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	latest   = -1
	testFlag = flag.Bool("t", false, "Whether or not to use test database")
	user     = -1
)

func metricsHandler() gin.HandlerFunc {
	reg := prometheus.NewRegistry()
	monitoring.NewMetrics(reg)

	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")

	router.GET("/metrics", metricsHandler())
	router.GET("/ping", controllers.Ping)
	router.GET("/", controllers.IndexTimeline)
	router.GET("/public", controllers.Timeline)
	router.GET("/:username", controllers.UserTimeline)
	router.GET("/:username/follow", controllers.UsernameFollow)
	router.GET("/:username/unfollow", controllers.UsernameUnfollow)
	router.POST("/register", controllers.Register_user)
	router.GET("/register", controllers.Register)
	router.POST("/add_message", controllers.AddMessage)
	router.POST("/login", controllers.Login_user)
	router.GET("/login", controllers.Loginf)
	router.GET("/logout", controllers.Logout_user)

	router.GET("/sim/latest", controllers.SimLatest)
	router.POST("/sim/register", controllers.SimRegister)
	router.GET("/sim/msgs", controllers.SimMsgs)
	router.POST("/sim/msgs/:username", controllers.SimPostUserMsg)
	router.GET("/sim/msgs/:username", controllers.SimGetUserMsg)
	router.POST("/sim/fllws/:username", controllers.SimPostUserFllws)
	router.GET("/sim/fllws/:username", controllers.SimGetUserFllws)

	return router

}

func main() {

	flag.Parse()

	if *testFlag {
		database.ConnectToTestDatabase()
	} else {
		database.ConnectToDatabase()
	}

	database.MigrateEntities()
	database.SeedDatabase()

	router := SetupRouter()
	router.Run() // port 8080
}
