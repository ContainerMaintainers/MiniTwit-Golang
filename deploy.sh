#!/usr/bin/env bash

docker pull $DOCKER_USERNAME/minitwit:latest 
nohup ./minitwit > out.log &