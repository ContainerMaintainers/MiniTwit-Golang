# MiniTwit-Golang

A Golang Gin refactor of the minitwit application.

## How to run the application

1. Clone the repo
2. Create an `.env` file from `.env.template` to set up the environment variables, docker won't be able to build without it.
3. Build project with docker. This will start the app and an instance of the database with example tweets. 

To build MiniTwit, you can use docker commands, or the shorcuts provided by the `Makefile`:
- `make` same as `docker compose up`. Build from docker-compose
- `make stop` same as `docker compose down`. Stops and removes containers, networks, volumes, and images created by.
- `make fresh_start` same as `docker compose up --build`. Rebuilds images and starts containers. 

Alternatively, if one has go installed, it is possible to run the application through: `go run minitwit.go`.
This approach assumes that a postgres database already exists. It is possible to run the application with an in memory database by using the -t flag: `go run minitwit.go -t`

By default, the application is run on port 8080.

## How to test the application

Using `go test` will run the `minitwit_test.go` tests. They test the traditional (not simulation-api) endpoints.

To test the simulation api, the application is is assumed to be running on port 8080 (see previous section).
Python has to be installed, as well as pytest (`pip install pytest`)
The tests are then run with `python -m pytest minitwit_sim_api_test.py`

To run the simulator: `python minitwit_simulator.py http://localhost:8080/sim`
