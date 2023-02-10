#!/bin/bash

docker build . -t minitwit-golang
docker run -i -t -p 8080:8080 minitwit-golang