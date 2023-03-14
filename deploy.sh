#!/usr/bin/env bash

docker pull $DOCKER_USERNAME/minitwit:latest
docker stop minitwit
docker run --rm -d -p $PORT:$PORT --name minitwit $DOCKER_USERNAME/minitwit:latest