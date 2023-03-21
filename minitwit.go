package main

import (
	"flag"
	"github.com/ContainerMaintainers/MiniTwit-Golang/controllers"
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"

	"net/http"
	// "github.com/prometheus/client_golang/prometheus"
    // "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	latest   = -1
	testFlag = flag.Bool("t", false, "Whether or not to use test database")
	user     = -1
)

func init() {
	initializers.LoadEnvVars()
}

func main() {

	http.Handle("/metrics", promhttp.Handler())

	flag.Parse()

	if *testFlag {
		database.ConnectToTestDatabase()
	} else {
		database.ConnectToDatabase()
	}

	database.MigrateEntities()
	database.SeedDatabase()

	router := controllers.SetupRouter()
	router.Run() // port 8080
}
