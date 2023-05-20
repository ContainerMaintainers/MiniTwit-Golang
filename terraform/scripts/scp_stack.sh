#!/bin/bash

# .env file
env='stack/minitwit_stack.yml'

# ssh key
key_file='ssh_key/terraform'

# swarm leader ip
ip=$(terraform output -raw minitwit-swarm-leader-ip-address)

# scp the file
scp -i $key_file $env root@$ip:/root/minitwit_stack.yml
