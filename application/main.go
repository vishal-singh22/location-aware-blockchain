// package main

// import (
//     "fmt"
//     "log"
//     "net/http"
//     "github.com/gorilla/mux"
// )

// func main() {
//     // Initialize the router
//     router := mux.NewRouter()

//     // Initialize blockchain connection
//     fabric, err := InitializeFabricClient()
//     if err != nil {
//         log.Fatalf("Failed to initialize Fabric client: %v", err)
//     }

//     // Create API handlers
//     handler := NewAPIHandler(fabric)

//     // Register routes
//     router.HandleFunc("/api/device", handler.RegisterDevice).Methods("POST")
//     router.HandleFunc("/api/device/{id}", handler.GetDevice).Methods("GET")
//     router.HandleFunc("/api/device/{id}/reputation", handler.UpdateReputation).Methods("PUT")
//     router.HandleFunc("/api/consensus/status", handler.GetConsensusStatus).Methods("GET")

//     // Start server
//     log.Printf("Starting server on :3000")
//     log.Fatal(http.ListenAndServe(":3000", router))
// }
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Device represents a device object.
type Device struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Reputation float64 `json:"reputation"`
}

// APIHandler holds dependencies for handling API requests.
type APIHandler struct {
	FabricClient *FabricClient
}

// FabricClient is a placeholder for the Hyperledger Fabric client.
type FabricClient struct {
	// Add fields for Fabric connection here (e.g., SDK instance).
}

// InitializeFabricClient initializes a connection to the Hyperledger Fabric network.
func InitializeFabricClient() (*FabricClient, error) {
	// Replace with actual Fabric SDK initialization logic.
	fmt.Println("Initializing Fabric client...")
	return &FabricClient{}, nil
}

// NewAPIHandler creates a new APIHandler instance.
func NewAPIHandler(client *FabricClient) *APIHandler {
	return &APIHandler{FabricClient: client}
}

// RegisterDevice handles the registration of a new device.
func (h *APIHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Logic to register the device in Hyperledger Fabric.
	fmt.Printf("Registering device: %+v\n", device)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Device registered successfully"})
}

// GetDevice retrieves a device by its ID.
func (h *APIHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]

	// Placeholder logic to fetch device information from Fabric.
	fmt.Printf("Fetching device with ID: %s\n", deviceID)
	device := Device{ID: deviceID, Name: "Sample Device", Reputation: 4.5}

	json.NewEncoder(w).Encode(device)
}

// UpdateReputation updates the reputation of a device.
func (h *APIHandler) UpdateReputation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]

	var update struct {
		Reputation float64 `json:"reputation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Placeholder logic to update device reputation in Fabric.
	fmt.Printf("Updating reputation for device %s to %f\n", deviceID, update.Reputation)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reputation updated successfully"})
}

// GetConsensusStatus retrieves the status of the consensus process.
func (h *APIHandler) GetConsensusStatus(w http.ResponseWriter, r *http.Request) {
	// Placeholder logic to fetch consensus status.
	fmt.Println("Fetching consensus status...")
	status := map[string]interface{}{
		"status":  "active",
		"shards":  5,
		"leader":  "node-123",
		"latency": "50ms",
	}

	json.NewEncoder(w).Encode(status)
}

func main() {
	// Initialize the router
	router := mux.NewRouter()

	// Initialize blockchain connection
	fabric, err := InitializeFabricClient()
	if err != nil {
		log.Fatalf("Failed to initialize Fabric client: %v", err)
	}

	// Create API handlers
	handler := NewAPIHandler(fabric)

	// Register routes
	router.HandleFunc("/api/device", handler.RegisterDevice).Methods("POST")
	router.HandleFunc("/api/device/{id}", handler.GetDevice).Methods("GET")
	router.HandleFunc("/api/device/{id}/reputation", handler.UpdateReputation).Methods("PUT")
	router.HandleFunc("/api/consensus/status", handler.GetConsensusStatus).Methods("GET")

	// Start server
	log.Printf("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
