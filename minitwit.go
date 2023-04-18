package main

import (
	"flag"
	"github.com/ContainerMaintainers/MiniTwit-Golang/controllers"
	"github.com/ContainerMaintainers/MiniTwit-Golang/database"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"
)

var (
	latest   = -1
	testFlag = flag.Bool("t", false, "Whether or not to use test database")
	user     = -1
)

func main() {

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
