#!/bin/bash

# grafana config
grafana='../../grafana'

# ssh key
key_file='ssh_key/terraform'

# ugly list concatenating of ips from terraform output
rows=$(terraform output -raw minitwit-swarm-leader-ip-address)
rows+=' '
rows+=$(terraform output -json minitwit-swarm-manager-ip-address | jq -r .[])
rows+=' '
rows+=$(terraform output -json minitwit-swarm-worker-ip-address | jq -r .[])

# scp the file
for ip in $rows; do
    scp -i $key_file $grafana root@$ip:/root/grafana
done