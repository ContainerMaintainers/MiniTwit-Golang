#!/usr/bin/sh

ENV_FILE=/.env
if [ -f "$ENV_FILE" ]; then
    echo ".env already exists"
else
    touch .env
    echo "DB_USER=$1" >> .env;
    echo "DB_PASSWORD=$2" >> .env;
    echo "DB_NAME=$3" >> .env;
    echo "DB_PORT=$4" >> .env;
    echo "DB_HOST=$5" >> .env;
    echo "PORT=$6" >> .env;
    echo "SESSION_KEY=$7" >> .env;
    echo "GIN_MODE=$8" >> .env;
fi