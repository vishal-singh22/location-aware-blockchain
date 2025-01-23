#!/bin/bash

# Set environment variables
export CHANNEL_NAME=mychannel
export CHAINCODE_NAME=iotchain
export CHAINCODE_VERSION=1.0
export CHAINCODE_SEQUENCE=1

# Package the chaincode
peer lifecycle chaincode package ${CHAINCODE_NAME}.tar.gz \
    --path ../chaincode \
    --lang golang \
    --label ${CHAINCODE_NAME}_${CHAINCODE_VERSION}

# Install the chaincode
peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz

# Approve the chaincode
peer lifecycle chaincode approveformyorg \
    -o localhost:7050 \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --version $CHAINCODE_VERSION \
    --sequence $CHAINCODE_SEQUENCE \
    --init-required \
    --package-id $CHAINCODE_NAME:$CHAINCODE_VERSION

# Commit the chaincode
peer lifecycle chaincode commit \
    -o localhost:7050 \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --version $CHAINCODE_VERSION \
    --sequence $CHAINCODE_SEQUENCE \
    --init-required