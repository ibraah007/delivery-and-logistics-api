package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type JobAssignment struct {
	JobID       string `json:"job_id"`
	Pickup      string `json:"pickup"`
	Destination string `json:"destination"`
	CargoType   string `json:"cargo_type"`
	AssignedAt  string `json:"assigned_at"`
	Status      string `json:"status"` // "ASSIGNED", "ACCEPTED", "IN_PROGRESS", "COMPLETED"
}

type DriverState struct {
	CarID          string        `json:"car_id"`
	DriverName     string        `json:"driver_name"`
	Role           string        `json:"role"` // "Long Haul", "Local Delivery", "Cold Chain"
	Status         string        `json:"status"`
	CurrentTask    string        `json:"current_task"`
	IncidentType   string        `json:"incident_type"`
	IncidentNotes  string        `json:"incident_notes"`
	EmergencyState string        `json:"emergency_state"` // "NONE", "ACTIVE_SOS", "MANAGER_ACKNOWLEDGED", "RESOLVED"
	ManagerNote    string        `json:"manager_note"`
	SafetyScore    int           `json:"safety_score"`
	PreTripPassed  bool          `json:"pre_trip_passed"`
	ActiveJob      *JobAssignment `json:"active_job"`
	Lat            float64       `json:"lat"`
	Lng            float64       `json:"lng"`
	Speed          float64       `json:"speed"`
	Heading        float64       `json:"heading"`
	LastUpdated    time.Time     `json:"last_updated"`
}

var (
	fleetData = make(map[string]*DriverState)
	mu        sync.RWMutex
)

// Driver updates telemetry & tasks
func updateFleetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload DriverState
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	existing, exists := fleetData[payload.CarID]
	if exists {
		// Preserve manager instructions if driver is just updating telemetry
		if payload.ManagerNote == "" {
			payload.ManagerNote = existing.ManagerNote
		}
		if payload.ActiveJob == nil {
			payload.ActiveJob = existing.ActiveJob
		}
	}
	payload.LastUpdated = time.Now()
	fleetData[payload.CarID] = &payload
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Manager assigns jobs, updates emergency response, or changes roles
func managerActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CarID          string         `json:"car_id"`
		Role           string         `json:"role"`
		Action         string         `json:"action"` // "ASSIGN_JOB", "ACK_SOS", "RESOLVE_SOS"
		ManagerNote    string         `json:"manager_note"`
		Job            *JobAssignment `json:"job"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	d, exists := fleetData[req.CarID]
	if !exists {
		d = &DriverState{CarID: req.CarID, DriverName: "Unassigned", SafetyScore: 100}
		fleetData[req.CarID] = d
	}

	if req.Role != "" {
		d.Role = req.Role
	}

	switch req.Action {
	case "ASSIGN_JOB":
		if req.Job != nil {
			req.Job.Status = "ASSIGNED"
			req.Job.AssignedAt = time.Now().Format("15:04:05")
			d.ActiveJob = req.Job
			d.ManagerNote = fmt.Sprintf("New Job Assigned: %s to %s", req.Job.Pickup, req.Job.Destination)
		}
	case "ACK_SOS":
		d.EmergencyState = "MANAGER_ACKNOWLEDGED"
		d.ManagerNote = req.ManagerNote
	case "RESOLVE_SOS":
		d.EmergencyState = "RESOLVED"
		d.Status = "ON_DUTY"
		d.IncidentType = "None"
		d.ManagerNote = "Incident marked resolved by Operations Head."
	}
	d.LastUpdated = time.Now()
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "action_processed"})
}

func getFleetHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleetData)
}

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/api/driver/update", updateFleetHandler)
	http.HandleFunc("/api/manager/action", managerActionHandler)
	http.HandleFunc("/api/fleet", getFleetHandler)

	fmt.Println("🚀 Mobikey Interactive Fleet Server running on http://localhost:8080")
	fmt.Println("👉 Manager Command: http://localhost:8080/dashboard.html")
	fmt.Println("👉 Driver Mobile App: http://localhost:8080/driver.html")
	http.ListenAndServe(":8080", nil)
}
