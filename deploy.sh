#!/usr/bin/env bash

echo "Pulling latest version"
docker pull $1/minitwit:latest
echo "Stopping current minitwit"
docker stop minitwit
echo "Deploying $DOCKER_USERNAME/minitwit:latest to $PORT"
docker run --rm -d -p $2:$2 $1/minitwit:latest
