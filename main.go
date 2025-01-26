package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WebhookPayload struct {
	Name string `json:"Name"`
	NotificationUsername string `json:"NotificationUsername"`
	UserId string `json:"UserId"`
	ServerName string `json:"ServerName"`
}

func main() {
	// Create a new instance of the WebhookPayload struct
	http.HandleFunc("/webhooks/jellyfin", func(w http.ResponseWriter, r *http.Request) {
		payload := WebhookPayload{}
		// Decode the JSON payload from the request body
		err := json.NewDecoder(r.Body).Decode(&payload)

		// log r.Body to see the payload
		if err != nil {
			fmt.Println("Error decoding JSON payload: ", err)
			return
		}
		// Print the decoded payload
		fmt.Printf("Received webhook payload: %+v\n", payload)

		// Respond with a 200 status code
		w.WriteHeader(http.StatusOK)
	})

	// Start the server on port 9090
	fmt.Println("Server listening on port 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}