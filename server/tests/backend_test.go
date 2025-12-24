package tests

import (
	"testing"

	"github.com/Zadigo/zanalytics/backend"
)

func TestNewPostgresDatabaseWithUrl(t *testing.T) {
	// Configuration with url
	var config = &backend.ServerBackendsConfig{
		Postgres: &backend.ServerBackendConfig{
			Url: "postgres://testuser:testpassword@localhost:5432",
		},
	}

	_, err := backend.NewPostgresDatabase(config)

	if err != nil {
		t.Fatalf("Failed to connect to Postgres database: %v", err)
	}
}

func TestCreateTables(t *testing.T) {
	backends := &backend.ServerBackendsConfig{
		Database: struct {
			Client string "json:\"client\" yaml:\"client\""
		}{Client: "postgres"},
		Postgres: &backend.ServerBackendConfig{
			Url: "postgres://testuser:testpassword@localhost:5432",
		},
	}

	conn, _ := backend.NewPostgresDatabase(backends)
	defer conn.Close(t.Context())

	serverConfig := &backend.ServerConfig{
		Config: struct {
			Endpoint     string "json:\"endpoint\" yaml:\"endpoint\""
			Port         int    "json:\"port\" yaml:\"port\""
			Username     string "json:\"username\" yaml:\"username\""
			Password     string "json:\"password\" yaml:\"password\""
			ClientId     string "json:\"client_id\" yaml:\"client_id\""
			ClientToken  string "json:\"client_token\" yaml:\"client_token\""
			MaxRetention int    "json:\"max_retention\" yaml:\"max_retention\""
			Country      string "json:\"country\" yaml:\"country\""
			Timezoze     string "json:\"timezone\" yaml:\"timezone\""
			Legal        struct {
				PolicyUrl  string "json:\"policy_url\" yaml:\"policy_url\""
				PrivacyUrl string "json:\"privacy_url\" yaml:\"privacy_url\""
			}
			Domains  []string                      "json:\"domains\" yaml:\"domains\""
			Backends *backend.ServerBackendsConfig "json:\"backends\" yaml:\"backends\""
		}{
			Username: "testuser",
			Password: "testuser",
		},
	}

	backend.CreateTables(conn, serverConfig)
}

func TestCreateUser(t *testing.T) {
	var config = &backend.ServerBackendsConfig{
		Postgres: &backend.ServerBackendConfig{
			Url: "postgres://testuser:testpassword@localhost:5432",
		},
	}

	conn, _ := backend.NewPostgresDatabase(config)
	defer conn.Close(t.Context())

	err := backend.CreateUser(conn, "testuser", "testpassword")

	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
}

func TestAuthenticateUser(t *testing.T) {
	var config = &backend.ServerBackendsConfig{
		Postgres: &backend.ServerBackendConfig{
			Url: "postgres://testuser:testpassword@localhost:5432",
		},
	}

	conn, _ := backend.NewPostgresDatabase(config)
	defer conn.Close(t.Context())

	err := backend.CreateUser(conn, "testuser", "testpassword")

	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	authenticated := backend.AuthenticateUser(conn, "testuser", "testpassword")

	if !authenticated {
		t.Fatalf("Failed to authenticate user")
	}
}

func TestNewRedisClient(t *testing.T) {
	var config = &backend.ServerBackendConfig{}

	client, err := backend.NewRedisClient(config)

	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}

	_, err = client.Ping(t.Context()).Result()
	if err != nil {
		t.Fatalf("Failed to ping Redis server: %v", err)
	}
}
