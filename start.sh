#!/bin/bash

# Build the Docker image
docker build -t parkpow_websocket .

# Run the Docker container with the .env file
docker run --env-file go_websocket/.env -p 8080:8080 parkpow_websocket
