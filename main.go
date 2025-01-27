package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

type WebhookPayload struct {
	Name string `json:"Name"`
	NotificationUsername string `json:"NotificationUsername"`
	UserId string `json:"UserId"`
	ServerName string `json:"ServerName"`
	NotificationType string `json:"NotificationType"`
	ItemId string `json:"ItemId"`
	ItemType string `json:"ItemType"`
	SeriesName string `json:"SeriesName"`
	SeasonNumber any `json:"SeasonNumber"`
	SeasonNumber00 any `json:"SeasonNumber00"`
	EpisodeNumber00 any `json:"EpisodeNumber00"`
	Thumbnail string `json:"Thumbnail"`
	AirTime string `json:"AirTime"`
}

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Println("Error creating Kafka producer: ", err)
		return
	}

	defer producer.Close()

	log.Println("Kafka producer created")

	http.HandleFunc("/webhooks/jellyfin", func(w http.ResponseWriter, r *http.Request) {
		payload := WebhookPayload{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println("Error decoding webhook payload: ", err)
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		payloadJSON, err := json.Marshal(payload)

		if err != nil {
			log.Println("Error decoding webhook payload: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		msq := &sarama.ProducerMessage{
			Topic: "jellyfin-notifications",
			Value: sarama.StringEncoder(payloadJSON),
		}

		partition, offset, err := producer.SendMessage(msq)

		if err != nil {
			log.Println("Error sending message to Kafka: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Server listening on port 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}