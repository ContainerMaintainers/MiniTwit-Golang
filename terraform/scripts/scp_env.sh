#!/bin/bash

# .env file
env='../.env'

# ssh key
key_file='ssh_key/terraform'

# ugly list concatenating of ips from terraform output
rows=$(terraform output -raw minitwit-swarm-leader-ip-address)
rows+=' '
rows+=$(terraform output -json minitwit-swarm-manager-ip-address | jq -r .[])
rows+=' '
rows+=$(terraform output -json minitwit-swarm-worker-ip-address | jq -r .[])

# rsync the file
for ip in $rows; do
    rsync -a -e "ssh -o 'StrictHostKeyChecking no' -i $key_file" $env root@$ip:/root/.env
done