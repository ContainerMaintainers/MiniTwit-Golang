#!/usr/bin/env bash

docker pull $DOCKER_USERNAME/minitwit:latest 
docker run --rm -d -p $PORT $DOCKER_USERNAME/minitwit > out.log &