package main

import (
    "encoding/json"
    "fmt"
    "time"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DeviceManager handles device-related operations
type DeviceManager struct {
    contractapi.Contract
}

// DeviceState represents the current state of a device
type DeviceState struct {
    ID              string    `json:"id"`
    Location        string    `json:"location"`
    ZoneID          string    `json:"zoneId"`
    Reputation      float64   `json:"reputation"`
    Status          string    `json:"status"` // "active", "inactive", "maintenance"
    LastUpdate      time.Time `json:"lastUpdate"`
    TransactionCount int      `json:"transactionCount"`
    SuccessfulTx    int      `json:"successfulTransactions"`
    FailedTx        int      `json:"failedTransactions"`
}

// DeviceTransaction represents a transaction performed by a device
type DeviceTransaction struct {
    DeviceID    string    `json:"deviceId"`
    Type        string    `json:"type"`
    Timestamp   time.Time `json:"timestamp"`
    Status      string    `json:"status"`
    ResponseTime int64    `json:"responseTime"` // in milliseconds
}

// CreateDevice initializes a new device in the system
func (dm *DeviceManager) CreateDevice(ctx contractapi.TransactionContextInterface, id string, location string, zoneId string) error {
    exists, err := dm.DeviceExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("device already exists: %s", id)
    }

    device := DeviceState{
        ID:              id,
        Location:        location,
        ZoneID:          zoneId,
        Reputation:      1.0, // Initial reputation
        Status:          "active",
        LastUpdate:      time.Now(),
        TransactionCount: 0,
        SuccessfulTx:    0,
        FailedTx:        0,
    }

    deviceJSON, err := json.Marshal(device)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, deviceJSON)
}

// DeviceExists checks if a device exists in the ledger
func (dm *DeviceManager) DeviceExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    deviceJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read device: %v", err)
    }
    return deviceJSON != nil, nil
}

// UpdateDeviceStatus updates the status of a device
func (dm *DeviceManager) UpdateDeviceStatus(ctx contractapi.TransactionContextInterface, id string, status string) error {
    device, err := dm.GetDevice(ctx, id)
    if err != nil {
        return err
    }

    device.Status = status
    device.LastUpdate = time.Now()

    deviceJSON, err := json.Marshal(device)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, deviceJSON)
}

// GetDevice retrieves device information
func (dm *DeviceManager) GetDevice(ctx contractapi.TransactionContextInterface, id string) (*DeviceState, error) {
    deviceJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read device: %v", err)
    }
    if deviceJSON == nil {
        return nil, fmt.Errorf("device does not exist: %s", id)
    }

    var device DeviceState
    err = json.Unmarshal(deviceJSON, &device)
    if err != nil {
        return nil, err
    }

    return &device, nil
}

// RecordTransaction records a transaction performed by a device
func (dm *DeviceManager) RecordTransaction(ctx contractapi.TransactionContextInterface, id string, txType string, status string, responseTime int64) error {
    device, err := dm.GetDevice(ctx, id)
    if err != nil {
        return err
    }

    // Update transaction counts
    device.TransactionCount++
    if status == "success" {
        device.SuccessfulTx++
    } else {
        device.FailedTx++
    }

    // Calculate new reputation based on transaction success rate
    successRate := float64(device.SuccessfulTx) / float64(device.TransactionCount)
    responseTimeScore := 1.0 - (float64(responseTime) / 5000.0) // Assuming 5000ms as max acceptable time
    if responseTimeScore < 0 {
        responseTimeScore = 0
    }

    // New reputation calculation (70% success rate, 30% response time)
    device.Reputation = (successRate * 0.7) + (responseTimeScore * 0.3)
    device.LastUpdate = time.Now()

    deviceJSON, err := json.Marshal(device)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, deviceJSON)
}

// QueryDevicesByZone gets all devices in a specific zone
func (dm *DeviceManager) QueryDevicesByZone(ctx contractapi.TransactionContextInterface, zoneId string) ([]*DeviceState, error) {
    queryString := fmt.Sprintf(`{"selector":{"zoneId":"%s"}}`, zoneId)
    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var devices []*DeviceState
    for resultsIterator.HasNext() {
        queryResult, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var device DeviceState
        err = json.Unmarshal(queryResult.Value, &device)
        if err != nil {
            return nil, err
        }
        devices = append(devices, &device)
    }

    return devices, nil
}

// GetDeviceHistory retrieves the history of a device
func (dm *DeviceManager) GetDeviceHistory(ctx contractapi.TransactionContextInterface, id string) ([]DeviceState, error) {
    historyIterator, err := ctx.GetStub().GetHistoryForKey(id)
    if err != nil {
        return nil, err
    }
    defer historyIterator.Close()

    var records []DeviceState
    for historyIterator.HasNext() {
        modification, err := historyIterator.Next()
        if err != nil {
            return nil, err
        }

        var device DeviceState
        err = json.Unmarshal(modification.Value, &device)
        if err != nil {
            return nil, err
        }

        records = append(records, device)
    }

    return records, nil
}