#!/bin/bash

# Run the Docker container in detached mode
docker run -d -p 8080:8080 --env-file .env --name parkpow_container parkpow_app

echo "Docker container started!"
