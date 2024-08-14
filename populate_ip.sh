#!/bin/bash

# Get the Docker IP address
DOCKER_IP=$(docker network inspect bridge --format '{{range .IPAM.Config}}{{.Gateway}}{{end}}')

# Check if Docker IP was retrieved successfully
if [ -z "$DOCKER_IP" ]; then
  echo "Failed to retrieve Docker IP address."
  exit 1
fi

# Assign the Docker IP to the SERVER_IPADDR key in the .env file
ENV_FILES="./client/.env ./client-swarm/.env"

for env_file in ${ENV_FILES}; do
    export SERVER_IPADDR="$DOCKER_IP"
done

