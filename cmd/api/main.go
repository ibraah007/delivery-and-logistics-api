package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Driver represents a complete transport system staff asset
type Driver struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"` // Long-Haul, Last-Mile, Dispatcher, Fleet Supervisor, Loader
	AssignedTask string    `json:"assigned_task"`
	Status       string    `json:"status"` // Active, Idle, Off-Duty
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var (
	drivers = []Driver{
		{ID: 1, Name: "Ibrah Samwel", Role: "Fleet Supervisor", AssignedTask: "Overseeing Kisumu Hub Routes", Status: "Active", Latitude: -0.0917, Longitude: 34.7680, UpdatedAt: time.Now()},
		{ID: 2, Name: "Karis", Role: "Long-Haul Trucker", AssignedTask: "Transit: Nairobi to Malaba Border", Status: "Active", Latitude: -0.0932, Longitude: 34.7695, UpdatedAt: time.Now()},
		{ID: 3, Name: "John Doe", Role: "Last-Mile Delivery", AssignedTask: "Pending Package Drop-offs", Status: "Idle", Latitude: -0.0950, Longitude: 34.7700, UpdatedAt: time.Now()},
	}
	mu     sync.Mutex
	nextID = 4
)

func main() {
	router := http.NewServeMux()

	// CORS Middleware
	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// GET & POST: Base Drivers Route
	router.HandleFunc("/api/drivers", cors(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(drivers)

		case http.MethodPost:
			var d Driver
			if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			d.ID = nextID
			nextID++
			d.UpdatedAt = time.Now()
			if d.Status == "" {
				d.Status = "Idle"
			}
			drivers = append(drivers, d)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(d)
		}
	}))

	// PUT & DELETE: Specific Driver Actions
	router.HandleFunc("/api/drivers/action", cors(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		switch r.Method {
		case http.MethodPut: // UPDATE Role & Task Assignment
			var updates Driver
			if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			for i, d := range drivers {
				if d.ID == id {
					drivers[i].Role = updates.Role
					drivers[i].AssignedTask = updates.AssignedTask
					drivers[i].Status = updates.Status
					drivers[i].UpdatedAt = time.Now()
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(drivers[i])
					return
				}
			}
			http.Error(w, "Driver not found", http.StatusNotFound)

		case http.MethodDelete: // DELETE Driver from Registry
			for i, d := range drivers {
				if d.ID == id {
					drivers = append(drivers[:i], drivers[i+1:]...)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"status":"deleted"}`))
					return
				}
			}
			http.Error(w, "Driver not found", http.StatusNotFound)
		}
	}))

	log.Println("Starting Logistics Transport Engine...")
	log.Println("Live API endpoint ready at http://127.0.0.1:8080/api/drivers 🚀")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server connection crashed: %v", err)
	}
}
