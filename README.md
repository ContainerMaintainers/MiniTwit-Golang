# MiniTwit-Golang

A Golang Gin refactor of the minitwit application.

## Before running the application

Make sure docker is installed as well as docker compose.
Furthermore, the docker drive for loki is also required. It can be installed via:
`docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions`

## How to run the application

1. Clone the repo
2. Create an `.env` file from `.env.template` to set up the environment variables, docker won't be able to build without it.
2. Build project with docker. This will start the app and an instance of the database with example tweets.

To build the entire application, one has to use the command `docker-compose up`. It must be run with the following environment variables: DB_USER, DB_PASSWORD, DB_HOST, DB_NAME, DB_PORT, PORT, SESSION_KEY, GIN_MODE

The whole command would be as follows (the environment variables should be replaced with the actual values):
`DB_USER=<user> DB_PASSWORD=<password> DB_HOST=<host> DB_NAME=<name> DB_PORT=<dbport> PORT=<port> SESSION_KEY=<sessionkey> GIN_MODE=<ginmode> docker-compose up`

This can be done more quickly by using shorcuts provided by the `Makefile`:
- `make` same as `docker compose up`. Build from docker-compose.
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

## Monitoring with Grafana and Prometheus
Prometheus can be accessed at [http://localhost:9090](http://localhost:9090)  
Grafana can be accessed at [http://localhost:3000](http://localhost:3000)  
Metrics end point at [http://localhost:8080/metrics](http://localhost:8080/metrics)