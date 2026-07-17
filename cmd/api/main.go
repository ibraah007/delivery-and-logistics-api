package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)

type Driver struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Role         string  `json:"role"`
	AssignedTask string  `json:"assigned_task"`
	Status       string  `json:"status"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

type DispatchPayload struct {
	Recipient string `json:"recipient"`
	Channel   string `json:"channel"`
	Contact   string `json:"contact"` 
	Customer  string `json:"customer"`
	Body      string `json:"body"`
}

var drivers = []Driver{
	{ID: 1, Name: "Ibrah Samwel", Role: "Fleet Supervisor", AssignedTask: "Overseeing Kisumu Hub Routes", Status: "Active", Latitude: -0.0917, Longitude: 34.7680},
	{ID: 2, Name: "Karis", Role: "Long-Haul Trucker", AssignedTask: "Transit: Nairobi to Malaba Border", Status: "Active", Latitude: -0.0932, Longitude: 34.7695},
}

func main() {
	http.HandleFunc("/api/drivers", driversHandler)
	http.HandleFunc("/api/dispatch", dispatchHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dashboard.html")
	})

	fmt.Println("Backend & UI System online running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func driversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func dispatchHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("[DISPATCH EVENT] Channel: %s | Target: %s (%s) | Job/Customer: %s\n", 
		payload.Channel, payload.Recipient, payload.Contact, payload.Customer)

	fromEmail := "your-email@gmail.com"       
	password := "your-app-password"           
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	atUsername := "sandbox"                  
	atAPIKey := "YOUR_AFRICAS_TALKING_KEY"    
	isSandbox := true                         

	if payload.Channel == "Email" {
		auth := smtp.PlainAuth("", fromEmail, password, smtpHost)
		messageContent := []byte(fmt.Sprintf("To: %s\r\n"+
			"Subject: Fleet Command Hub - Operational Dispatch Notification\r\n\r\n"+
			"Attention %s,\n\n"+
			"You have a new logistics task alert assigned.\n\n"+
			"Details:\n"+
			"- Associated Context: %s\n"+
			"- Task Manifest: %s\n", 
			payload.Contact, payload.Recipient, payload.Customer, payload.Body))

		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{payload.Contact}, messageContent)
		if err != nil {
			log.Println("SMTP Error:", err)
			http.Error(w, "Failed to connect and deliver message via mail server.", http.StatusInternalServerError)
			return
		}
	} else if payload.Channel == "SMS" {
		toPhone := payload.Contact
		if strings.HasPrefix(toPhone, "0") {
			toPhone = "+254" + toPhone[1:]
		}

		smsBody := fmt.Sprintf("Alert to %s: %s (Context: %s)", payload.Recipient, payload.Body, payload.Customer)

		apiURL := "https://api.africastalking.com/version1/messaging"
		if isSandbox {
			apiURL = "https://api.sandbox.africastalking.com/version1/messaging"
		}

		data := url.Values{}
		data.Set("username", atUsername)
		data.Set("to", toPhone)
		data.Set("message", smsBody)

		req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
		if err != nil {
			log.Println("SMS Request Prep Error:", err)
			http.Error(w, "Failed to prep SMS gateway request", http.StatusInternalServerError)
			return
		}

		req.Header.Add("apiKey", atAPIKey)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("SMS Gateway Delivery Error:", err)
			http.Error(w, "SMS Gateway offline", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			log.Printf("SMS Gateway rejected request with status: %d\n", resp.StatusCode)
			http.Error(w, "Gateway credentials invalid or rejected by telecom host.", http.StatusInternalServerError)
			return
		}
		log.Printf("[LIVE SMS SENT SUCCESSFULLY] Routed real message to: %s\n", toPhone)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Logistics movement transmission completed cleanly."}`))
}
