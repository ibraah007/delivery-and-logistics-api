package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
)

// Driver represents our worker asset
type Driver struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Role         string  `json:"role"`
	AssignedTask string  `json:"assigned_task"`
	Status       string  `json:"status"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

// DispatchPayload captures the frontend messaging form
type DispatchPayload struct {
	Recipient string `json:"recipient"`
	Channel   string `json:"channel"`
	Customer  string `json:"customer"`
	Body      string `json:"body"`
}

var drivers = []Driver{
	{ID: 1, Name: "Ibrah Samwel", Role: "Fleet Supervisor", AssignedTask: "Overseeing Kisumu Hub Routes", Status: "Active", Latitude: -0.0917, Longitude: 34.7680},
}

func main() {
	http.HandleFunc("/api/drivers", driversHandler)
	http.HandleFunc("/api/dispatch", dispatchHandler) // New Route for actual alerts

	fmt.Println("Backend system online running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func driversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func dispatchHandler(w http.ResponseWriter, r *http.Request) {
	// Setup CORS headers so your frontend dashboard can talk to it cleanly
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload DispatchPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ----------------------------------------------------------------
	// REAL EMAIL SMTP SETUP
	// ----------------------------------------------------------------
	fromEmail := "your-email@gmail.com" // Change to your Gmail address
	password := "your-app-password"     // Change to your 16-character Google App Password
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Hardcoded destination mail targets for your test
	driverEmail := "your-driver-email@example.com"   // Put your actual email address here
	customerEmail := "your-friend-email@example.com" // Put your friend's email address here

	auth := smtp.PlainAuth("", fromEmail, password, smtpHost)

	// Formulating clear, clean body configurations for both target users
	driverMessage := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Fleet Matrix System - New Job Dispatch Alert\r\n\r\n"+
		"Hello %s,\n\nYou have been assigned a new logistics movement task.\n"+
		"Job Context / Customer Reference: %s\n\n"+
		"Operational Instructions:\n%s\n\n"+
		"Please check your active fleet console map immediately to begin route logging.\n",
		driverEmail, payload.Recipient, payload.Customer, payload.Body))

	customerMessage := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Logistics Notification - Order Movement Active\r\n\r\n"+
		"Hello %s,\n\nGood news! Your delivery agent (%s) has started handling your asset route.\n"+
		"Task details assigned by Logistics Manager: \"%s\"\r\n",
		customerEmail, payload.Customer, payload.Recipient, payload.Body))

	// Execute real email dispatch to Driver
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{driverEmail}, driverMessage)
	if err != nil {
		log.Println("Driver email transmission failed:", err)
		http.Error(w, "Failed to send notification to driver", http.StatusInternalServerError)
		return
	}

	// Execute real email dispatch to Customer
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{customerEmail}, customerMessage)
	if err != nil {
		log.Println("Customer email transmission failed:", err)
		http.Error(w, "Failed to send notification to customer", http.StatusInternalServerError)
		return
	}

	// Return clean verification back to our frontend terminal log
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Real-world dual email dispatches delivered successfully!"}`))
}
