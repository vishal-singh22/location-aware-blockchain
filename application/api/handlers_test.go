package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockFabricClient struct {
    mock.Mock
}

func (m *MockFabricClient) RegisterDevice(id, location, zoneId string) error {
    args := m.Called(id, location, zoneId)
    return args.Error(0)
}

func (m *MockFabricClient) QueryDevice(id string) (*DeviceState, error) {
    args := m.Called(id)
    return args.Get(0).(*DeviceState), args.Error(1)
}

func (m *MockFabricClient) UpdateDeviceReputation(id string, reputation float64) error {
    args := m.Called(id, reputation)
    return args.Error(0)
}

func TestAPIHandlers(t *testing.T) {
    t.Run("RegisterDevice", func(t *testing.T) {
        // Test successful registration
        t.Run("Success", func(t *testing.T) {
            mockClient := new(MockFabricClient)
            handler := NewAPIHandler(mockClient)

            mockClient.On("RegisterDevice", "test-device", "test-zone", "Z1").Return(nil)

            body := DeviceRegistration{
                ID: "test-device",
                Location: "test-zone",
                ZoneID: "Z1",
            }
            bodyJSON, _ := json.Marshal(body)

            req := httptest.NewRequest("POST", "/api/device", bytes.NewBuffer(bodyJSON))
            rec := httptest.NewRecorder()

            handler.RegisterDevice(rec, req)

            assert.Equal(t, http.StatusCreated, rec.Code)
            mockClient.AssertExpectations(t)
        })

        // Test invalid request body
        t.Run("InvalidBody", func(t *testing.T) {
            mockClient := new(MockFabricClient)
            handler := NewAPIHandler(mockClient)

            req := httptest.NewRequest("POST", "/api/device", bytes.NewBuffer([]byte("invalid json")))
            rec := httptest.NewRecorder()

            handler.RegisterDevice(rec, req)

            assert.Equal(t, http.StatusBadRequest, rec.Code)
        })
    })

    t.Run("GetDevice", func(t *testing.T) {
        // Test successful device retrieval
        t.Run("Success", func(t *testing.T) {
            mockClient := new(MockFabricClient)
            handler := NewAPIHandler(mockClient)

            device := &DeviceState{
                ID: "test-device",
                Location: "test-zone",
                Reputation: 1.0,
            }
            mockClient.On("QueryDevice", "test-device").Return(device, nil)

            req := httptest.NewRequest("GET", "/api/device/test-device", nil)
            rec := httptest.NewRecorder()

            vars := map[string]string{
                "id": "test-device",
            }
            req = mux.SetURLVars(req, vars)

            handler.GetDevice(rec, req)

            assert.Equal(t, http.StatusOK, rec.Code)
            mockClient.AssertExpectations(t)
        })
    })

    t.Run("UpdateReputation", func(t *testing.T) {
        // Test successful reputation update
        t.Run("Success", func(t *testing.T) {
            mockClient := new(MockFabricClient)
            handler := NewAPIHandler(mockClient)

            mockClient.On("UpdateDeviceReputation", "test-device", 0.9).Return(nil)

            body := ReputationUpdate{
                NewReputation: 0.9,
            }
            bodyJSON, _ := json.Marshal(body)

            req := httptest.NewRequest("PUT", "/api/device/test-device/reputation", bytes.NewBuffer(bodyJSON))
            rec := httptest.NewRecorder()

            vars := map[string]string{
                "id": "test-device",
            }
            req = mux.SetURLVars(req, vars)

            handler.UpdateReputation(rec, req)

            assert.Equal(t, http.StatusOK, rec.Code)
            mockClient.AssertExpectations(t)
        })
    })
}