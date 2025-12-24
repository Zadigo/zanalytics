package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Zadigo/zanalytics/backend"
	"github.com/gorilla/websocket"
	yamlParser "gopkg.in/yaml.v3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handles websocket connections for live analytics data
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

// Handles incoming analytics data via HTTP POST requests
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(body))
}

func beforeStart() (*backend.ServerConfig, error) {
	log.Println("Performing pre-startup tasks...")

	// Load configuration from YAML file
	buffer, err := os.ReadFile("config.yaml")

	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
		return nil, err
	}

	content := &backend.ServerConfig{}
	err = yamlParser.Unmarshal(buffer, content)

	if err != nil {
		log.Fatalf("Error parsing config file: %v\n", err)
		return nil, err
	}

	// Connect to the prefered database (Postgres, SQLite, etc.)
	conn, err := backend.NewPostgresDatabase(content.Config.Backends)

	if err != nil {
		log.Fatalf("Failed to connect to Postgres database: %v\n", err)
	}

	defer conn.Close(context.Background())
	err = conn.Ping(context.Background())

	if err != nil {
		log.Fatalf("Failed to ping Postgres database: %v\n", err)
		return nil, err
	}

	// Connect to RabbitMQ and start the consumer server

	// Connect to Redis

	log.Println("☑ Pre-startup tasks completed.")

	return content, nil
}

// Entry point of the Analytics Service
func main() {
	log.Println("🚀 Starting Analytics Service on port 9000...")
	log.Println("🫆 HTTP Endpoint: http://127.0.0.1:9000/v1/analytics")
	log.Println("🫆 Websocket Endpoint: ws://127.0.0.1:9000/v1/analytics/live")

	// Setting up HTTP handlers
	http.HandleFunc("/v1/analytics", analyticsHandler)
	http.HandleFunc("/v1/analytics/live", liveAnalyticsHandler)

	// Perform pre-startup tasks
	go beforeStart()

	// Starting the HTTP server
	err := http.ListenAndServe(":9000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
