package main

import (
    "encoding/json"
	"fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)
 
// SmartContract provides functions for managing IoT devices
type SmartContract struct {
    contractapi.Contract
}

// Device represents an IoT device with location and reputation
type Device struct {
    ID         string  `json:"id"`
    Location   string  `json:"location"`
    Reputation float64 `json:"reputation"`
    LastUpdate int64   `json:"lastUpdate"`
    ZoneID     string  `json:"zoneId"`
}

// InitLedger adds a base set of devices to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    devices := []Device{
        {
            ID:         "device1",
            Location:   "zone1",
            Reputation: 1.0,
            LastUpdate: 1635724800,
            ZoneID:     "Z1",
        },
    }

    for _, device := range devices {
        deviceJSON, err := json.Marshal(device)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(device.ID, deviceJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state: %v", err)
        }
    }

    return nil
}

// RegisterDevice adds a new device to the world state
func (s *SmartContract) RegisterDevice(ctx contractapi.TransactionContextInterface, id string, location string, zoneId string) error {
    device := Device{
        ID:         id,
        Location:   location,
        Reputation: 1.0, // Initial reputation
        LastUpdate: ctx.GetStub().GetTxTimestamp().Seconds,
        ZoneID:     zoneId,
    }

    deviceJSON, err := json.Marshal(device)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, deviceJSON)
}

// QueryDevice returns the device stored in the world state with given id
func (s *SmartContract) QueryDevice(ctx contractapi.TransactionContextInterface, id string) (*Device, error) {
    deviceJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if deviceJSON == nil {
        return nil, fmt.Errorf("the device %s does not exist", id)
    }

    var device Device
    err = json.Unmarshal(deviceJSON, &device)
    if err != nil {
        return nil, err
    }

    return &device, nil
}

// UpdateDeviceReputation updates the reputation of a device
func (s *SmartContract) UpdateDeviceReputation(ctx contractapi.TransactionContextInterface, id string, newReputation float64) error {
    device, err := s.QueryDevice(ctx, id)
    if err != nil {
        return err
    }

    device.Reputation = newReputation
    device.LastUpdate = ctx.GetStub().GetTxTimestamp().Seconds

    deviceJSON, err := json.Marshal(device)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, deviceJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        fmt.Printf("Error creating chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %s", err.Error())
    }
}