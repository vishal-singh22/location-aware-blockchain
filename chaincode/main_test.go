package main

import (
    "encoding/json"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MockStub implements the chaincode stub interface
type MockStub struct {
    mock.Mock
    shim.ChaincodeStubInterface
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
    args := ms.Called(key)
    return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
    args := ms.Called(key, value)
    return args.Error(0)
}

func (ms *MockStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
    return &timestamp.Timestamp{Seconds: time.Now().Unix()}, nil
}

// Test cases
func TestSmartContract(t *testing.T) {
    t.Run("RegisterDevice", func(t *testing.T) {
        // Test successful device registration
        t.Run("Success", func(t *testing.T) {
            contract := new(SmartContract)
            ctx := new(MockStub)

            device := Device{
                ID: "test-device",
                Location: "test-zone",
                Reputation: 1.0,
                ZoneID: "Z1",
            }

            deviceJSON, _ := json.Marshal(device)
            ctx.On("GetState", "test-device").Return([]byte{}, nil)
            ctx.On("PutState", "test-device", deviceJSON).Return(nil)

            err := contract.RegisterDevice(ctx, "test-device", "test-zone", "Z1")
            assert.NoError(t, err)
            ctx.AssertExpectations(t)
        })

        // Test duplicate device registration
        t.Run("DuplicateDevice", func(t *testing.T) {
            contract := new(SmartContract)
            ctx := new(MockStub)

            existingDevice := Device{
                ID: "test-device",
                Location: "test-zone",
                Reputation: 1.0,
            }
            deviceJSON, _ := json.Marshal(existingDevice)

            ctx.On("GetState", "test-device").Return(deviceJSON, nil)

            err := contract.RegisterDevice(ctx, "test-device", "test-zone", "Z1")
            assert.Error(t, err)
            assert.Contains(t, err.Error(), "already exists")
        })
    })

    t.Run("QueryDevice", func(t *testing.T) {
        // Test successful device query
        t.Run("Success", func(t *testing.T) {
            contract := new(SmartContract)
            ctx := new(MockStub)

            device := Device{
                ID: "test-device",
                Location: "test-zone",
                Reputation: 1.0,
            }
            deviceJSON, _ := json.Marshal(device)

            ctx.On("GetState", "test-device").Return(deviceJSON, nil)

            result, err := contract.QueryDevice(ctx, "test-device")
            assert.NoError(t, err)
            assert.Equal(t, "test-device", result.ID)
            assert.Equal(t, "test-zone", result.Location)
        })

        // Test query non-existent device
        t.Run("DeviceNotFound", func(t *testing.T) {
            contract := new(SmartContract)
            ctx := new(MockStub)

            ctx.On("GetState", "non-existent").Return([]byte{}, nil)

            result, err := contract.QueryDevice(ctx, "non-existent")
            assert.Error(t, err)
            assert.Nil(t, result)
        })
    })

    t.Run("UpdateDeviceReputation", func(t *testing.T) {
        // Test successful reputation update
        t.Run("Success", func(t *testing.T) {
            contract := new(SmartContract)
            ctx := new(MockStub)

            device := Device{
                ID: "test-device",
                Location: "test-zone",
                Reputation: 1.0,
            }
            deviceJSON, _ := json.Marshal(device)

            updatedDevice := device
            updatedDevice.Reputation = 0.9
            updatedDeviceJSON, _ := json.Marshal(updatedDevice)

            ctx.On("GetState", "test-device").Return(deviceJSON, nil)
            ctx.On("PutState", "test-device", updatedDeviceJSON).Return(nil)

            err := contract.UpdateDeviceReputation(ctx, "test-device", 0.9)
            assert.NoError(t, err)
        })

        // Test invalid reputation value
        t.Run("InvalidReputation", func(t *testing.T) {
            contract := new(SmartContract)
            ctx := new(MockStub)

            err := contract.UpdateDeviceReputation(ctx, "test-device", 1.5)
            assert.Error(t, err)
        })
    })
}