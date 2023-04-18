#!/usr/bin/env bash

echo "Pulling latest version"
docker pull $1/minitwit:latest
echo "Stopping current minitwit"
docker stop minitwit
echo "Deploying $DOCKER_USERNAME/minitwit:latest to $PORT"
docker run -d -p $2:$2 --name minitwit $1/minitwit:latest --env DB_NAME=$3 --env DB_USER=$6 --env DB_PASSWORD=$5 --env DB_HOST=$7 --env DB_PORT=$4 --env PORT=$2
