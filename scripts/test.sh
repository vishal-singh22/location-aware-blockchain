#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Start the network
echo "Starting the network..."
./network.sh up

# Wait for network to stabilize
sleep 5

# Test API endpoints
test_api() {
    local endpoint=$1
    local method=$2
    local data=$3
    local expected_status=$4

    response=$(curl -s -X $method \
        -H "Content-Type: application/json" \
        -d "$data" \
        -w "%{http_code}" \
        http://localhost:3000$endpoint)

    status_code=${response: -3}
    if [ $status_code -eq $expected_status ]; then
        echo -e "${GREEN}✓ $method $endpoint - Success${NC}"
    else
        echo -e "${RED}✗ $method $endpoint - Failed${NC}"
        exit 1
    fi
}

# Run tests
echo "Testing API endpoints..."

# Test device registration
test_api "/api/device" "POST" '{"id":"test-device-1","location":"zone1","zoneId":"Z1"}' 201

# Test device query
test_api "/api/device/test-device-1" "GET" "" 200

# Test reputation update
test_api "/api/device/test-device-1/reputation" "PUT" '{"newReputation":0.95}' 200
!/bin/bash

