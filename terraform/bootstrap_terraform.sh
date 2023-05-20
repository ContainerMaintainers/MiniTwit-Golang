#!/bin/bash

echo -e "\n--> Bootstrapping Minitwit\n"

echo -e "\n--> Loading environment variables from secrets file\n"
source secrets

echo -e "\n--> Checking that environment variables are set\n"
# check that all variables are set
[ -z "$TF_VAR_do_token" ] && echo "TF_VAR_do_token is not set" && exit
[ -z "$SPACE_NAME" ] && echo "SPACE_NAME is not set" && exit
[ -z "$STATE_FILE" ] && echo "STATE_FILE is not set" && exit
[ -z "$AWS_ACCESS_KEY_ID" ] && echo "AWS_ACCESS_KEY_ID is not set" && exit
[ -z "$AWS_SECRET_ACCESS_KEY" ] && echo "AWS_SECRET_ACCESS_KEY is not set" && exit

echo -e "\n--> Initializing terraform\n"
# initialize terraform
terraform init \
    -backend-config "bucket=$SPACE_NAME" \
    -backend-config "key=$STATE_FILE" \
    -backend-config "access_key=$AWS_ACCESS_KEY_ID" \
    -backend-config "secret_key=$AWS_SECRET_ACCESS_KEY"

# check that everything looks good
echo -e "\n--> Validating terraform configuration\n"
terraform validate

# create infrastructure
echo -e "\n--> Creating Infrastructure\n"
terraform apply -auto-approve

# generate loadbalancer configuration
echo -e "\n--> Generating loadbalancer configuration\n"
bash scripts/gen_load_balancer_config.sh

# scp loadbalancer config to all nodes
echo -e "\n--> Copying loadbalancer configuration to nodes\n"
bash scripts/scp_load_balancer_config.sh

# .env file to all nodes
echo -e "\n--> Copying .env file to nodes\n"
bash scripts/scp_env.sh

# grafana config to all nodes
echo -e "\n--> Copying grafan config files to nodes\n"
bash scripts/scp_grafana_config.sh

# stack file to all nodes
echo -e "\n--> Copying stack file to leader\n"
bash scripts/scp_stack.sh

# prometheus config file to all nodes
echo -e "\n--> Copying prometheus config file to nodes\n"
bash scripts/scp_prometheus_config.sh

# deploy the stack to the cluster
echo -e "\n--> Deploying the Minitwit stack to the cluster\n"
ssh \
    -o 'StrictHostKeyChecking no' \
    root@$(terraform output -raw minitwit-swarm-leader-ip-address) \
    -i ssh_key/terraform \
    "DB_NAME=$DB_NAME DB_USER=$DB_USER DB_PASSWORD=$DB_PASSWORD DB_PORT=$DB_PORT SESSION_KEY=$SESSION_KEY GIN_MODE=${GIN_MODE} docker stack deploy minitwit -c minitwit_stack.yml"

echo -e "\n--> Done bootstrapping Minitwit"
echo -e "--> The dbs will need a moment to initialize, this can take up to a couple of minutes..."
echo -e "--> Site will be avilable @ http://$(terraform output -raw public_ip)"
echo -e "--> You can check the status of swarm cluster @ http://$(terraform output -raw minitwit-swarm-leader-ip-address):8888"
echo -e "--> ssh to swarm leader with 'ssh root@\$(terraform output -raw minitwit-swarm-leader-ip-address) -i ssh_key/terraform'"
echo -e "--> To remove the infrastructure run: terraform destroy -auto-approve"