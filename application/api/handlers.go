package api

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type APIHandler struct {
    fabricClient *FabricClient
}

type DeviceRegistration struct {
    ID       string `json:"id"`
    Location string `json:"location"`
    ZoneID   string `json:"zoneId"`
}

type ReputationUpdate struct {
    NewReputation float64 `json:"newReputation"`
}

func NewAPIHandler(fabric *FabricClient) *APIHandler {
    return &APIHandler{
        fabricClient: fabric,
    }
}

func (h *APIHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
    var reg DeviceRegistration
    if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err := h.fabricClient.RegisterDevice(reg.ID, reg.Location, reg.ZoneID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deviceID := vars["id"]

    device, err := h.fabricClient.QueryDevice(deviceID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(device)
}

func (h *APIHandler) UpdateReputation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deviceID := vars["id"]

    var update ReputationUpdate
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err := h.fabricClient.UpdateDeviceReputation(deviceID, update.NewReputation)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *APIHandler) GetConsensusStatus(w http.ResponseWriter, r *http.Request) {
    status, err := h.fabricClient.GetConsensusStatus()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(status)
}