#!/bin/bash

function networkDown() {
    echo "Stopping and removing containers..."
    docker stop orderer.example.com peer0.org1.example.com 2>/dev/null || true
    docker rm orderer.example.com peer0.org1.example.com 2>/dev/null || true
    docker network rm network_iot_network 2>/dev/null || true
    echo "Cleanup completed"
}

function networkUp() {
    echo "Starting network..."
    # First clean up any existing containers
    networkDown
    # Then start the network
    docker-compose -f ../network/docker-compose.yml up -d
}

case "$1" in
    "up")
        networkUp
        ;;
    "down")
        networkDown
        ;;
    *)
        echo "Usage: $0 [up|down]"
        exit 1
esac