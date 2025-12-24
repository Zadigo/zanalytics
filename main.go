package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func liveAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Could not start live analytics websocket:", err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("POST - Request received")

	// context := r.Context()
	// token := r.URL.Query().Get("token")

	// Only POST method is allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Only JSON content type is allowed
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; application/octet-stream; charset=utf-8")
	fmt.Fprint(w, string(body))
}

func main() {
	fmt.Println("🚀 Starting Analytics Service on port 9000...")
	fmt.Println("🫆 HTTP Endpoint: http://127.0.0.1:9000/v1/analytics")
	fmt.Println("🫆 Websocket Endpoint: ws://127.0.0.1:9000/v1/analytics/live")

	// Setting up HTTP handlers
	http.HandleFunc("/v1/analytics", analyticsHandler)
	http.HandleFunc("/v1/analytics/live", liveAnalyticsHandler)

	// Connect to the prefered database (Postgres, SQLite, etc.)

	// Connect to RabbitMQ and start the consumer server

	// Connect to Redis

	// Starting the HTTP server
	err := http.ListenAndServe(":9000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
