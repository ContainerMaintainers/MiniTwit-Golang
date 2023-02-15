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
