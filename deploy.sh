#!/usr/bin/env bash

# ARGUMENTS
# $1 : Docker hub username
# $2 : Application port 
# $3 : Database name
# $4 : Database port
# $5 : Database password
# $6 : Database user
# $7 : Database host
# $8 : Session key
# $9 : Gin mode
# $10 : Loki username
# $11 : Loki password
# $12 : Loki host
# $13 : Loki port

echo "Installing loki driver"
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions || true
echo "Pulling latest version"
docker pull $1/minitwit:latest
echo "Stopping current minitwit"
docker stop minitwit || true
docker rm minitwit || true
echo "Deploying $DOCKER_USERNAME/minitwit:latest to $PORT"
docker run -d \
	-p $2:$2 \
    --env PORT=$2 \
	--env DB_NAME=$3 \
    --env DB_PORT=$4 \
    --env DB_PASSWORD=$5 \
	--env DB_USER=$6 \
	--env DB_HOST=$7 \
    --env SESSION_KEY=$8 \
	--env GIN_MODE=$9 \
	--name minitwit \
	--log-driver=loki \
    --log-opt loki-url="http://${10}:${11}@${12}:${13}/loki/api/v1/push" \
    --log-opt loki-retries=5 \
    --log-opt loki-batch-size=400 \
    --name minitwit \
	$1/minitwit
