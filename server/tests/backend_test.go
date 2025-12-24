package tests

import (
	"testing"

	"github.com/Zadigo/zanalytics/backend"
)

func TestNewPostgresDatabaseWithUrl(t *testing.T) {
	// Configuration with url
	var config = &backend.ServerBackendConfig{
		Url: "postgres://testuser:testpassword@localhost:5432/zanalytics",
	}

	_, err := backend.NewPostgresDatabase(config)

	if err != nil {
		t.Fatalf("Failed to connect to Postgres database: %v", err)
	}
}

func TestNewPostgresDatabaseWithoutUrl(t *testing.T) {
	// Configuration without url
	var config = &backend.ServerBackendConfig{
		Username: "testuser",
		Password: "testuser",
	}

	conn, err := backend.NewPostgresDatabase(config)

	if err != nil {
		t.Fatalf("Failed to connect to Postgres database: %v", err)
	}

	if conn == nil {
		t.Fatalf("Connection is nil or not of type pointer")
	}
}

func TestCreateTables(t *testing.T) {
	var config = &backend.ServerBackendConfig{
		Url: "postgres://testuser:testpassword@localhost:5432/zanalytics",
	}

	conn, _ := backend.NewPostgresDatabase(config)
	defer conn.Close(t.Context())

	backend.CreateTables(conn, config)
}

func TestCreateUser(t *testing.T) {
	var config = &backend.ServerBackendConfig{
		Url: "postgres://testuser:testpassword@localhost:5432/zanalytics",
	}

	conn, _ := backend.NewPostgresDatabase(config)
	defer conn.Close(t.Context())

	err := backend.CreateUser(conn, "testuser", "testpassword")

	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
}

func TestAuthenticateUser(t *testing.T) {
	var config = &backend.ServerBackendConfig{
		Url: "postgres://testuser:testpassword@localhost:5432/zanalytics",
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
