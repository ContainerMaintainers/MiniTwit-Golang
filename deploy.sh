#!/usr/bin/env bash

echo "Pulling latest version"
docker pull $1/minitwit:latest
echo "Stopping current minitwit"
docker stop minitwit
echo "Deploying $DOCKER_USERNAME/minitwit:latest to $PORT"
docker run --rm -d -p $2:$2 --name minitwit $1/minitwit:latest -e DB_NAME=$3 -e DB_USER=$6 -e DB_PASSWORD=$5 -e DB_HOST=$7 -e DB_PORT=$4 -e PORT=$2
