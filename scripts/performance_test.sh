#!/bin/bash

# Number of concurrent requests
CONCURRENT=50
# Total number of requests
REQUESTS=1000

echo "Running performance tests..."

# Test device registration performance
ab -n $REQUESTS -c $CONCURRENT -T 'application/json' \
   -p device.json http://localhost:3000/api/device/

# Test query performance
ab -n $REQUESTS -c $CONCURRENT \
   http://localhost:3000/api/device/test-device-1