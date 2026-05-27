package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/models"
)

// TrackLocationHandler handles live GPS pings from drivers
func TrackLocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Extract the driver's user_id safely from the JWT context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized: Identity missing", http.StatusUnauthorized)
		return
	}

	// 2. Parse the latitude and longitude from the request body
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Latitude == 0 || req.Longitude == 0 {
		http.Error(w, "Invalid coordinates provided", http.StatusBadRequest)
		return
	}

	// 3. Build our Location model
	location := models.Location{
		DriverID:  userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: time.Now(),
	}

	// 4. Save the location ping to the PostgreSQL database using DBInstance
	query := `INSERT INTO locations (driver_id, latitude, longitude, timestamp) 
	          VALUES ($1, $2, $3, $4)`
	
	_, err = configs.DBInstance.Exec(query, location.DriverID, location.Latitude, location.Longitude, location.Timestamp)
	if err != nil {
		http.Error(w, "Failed to log location data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Send back a clean success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Live location telemetry logged successfully",
	})
}
