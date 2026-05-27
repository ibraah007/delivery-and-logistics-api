package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/models"
	"github.com/ibrah/logistics-tracking-api/internal/repository"
)

// TrackLocationHandler handles live GPS pings from drivers with strict validation
func TrackLocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized: Identity missing", http.StatusUnauthorized)
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Boundary Safeguard: Latitudes are [-90, 90], Longitudes are [-180, 180]
	if req.Latitude < -90.0 || req.Latitude > 90.0 || req.Longitude < -180.0 || req.Longitude > 180.0 {
		http.Error(w, "Bad Request: Coordinates out of realistic global boundaries", http.StatusBadRequest)
		return
	}

	location := models.Location{
		DriverID:  userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: time.Now(),
	}

	query := `INSERT INTO locations (driver_id, latitude, longitude, timestamp) 
	          VALUES ($1, $2, $3, $4)`

	_, err := configs.DBInstance.Exec(query, location.DriverID, location.Latitude, location.Longitude, location.Timestamp)
	if err != nil {
		http.Error(w, "Failed to log location data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Live location telemetry logged successfully",
	})
}

// GetTrackingHistoryHandler allows admins to review a driver's route telemetry safely
func GetTrackingHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	driverIDStr := r.URL.Query().Get("driver_id")
	if driverIDStr == "" {
		http.Error(w, "Missing required query parameter: driver_id", http.StatusBadRequest)
		return
	}

	var driverID uint
	_, err := fmt.Sscanf(driverIDStr, "%d", &driverID)
	if err != nil || driverID == 0 {
		http.Error(w, "Invalid driver_id format", http.StatusBadRequest)
		return
	}

	// Utilizing the updated timeout-secured repository function
	history, err := repository.GetDriverLocationHistorySecured(driverID)
	if err != nil {
		http.Error(w, "Failed to retrieve tracking data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if history == nil {
		history = []models.Location{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
