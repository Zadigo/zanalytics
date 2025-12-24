package backend

import (
	"context"
	"fmt"

	postGres "github.com/jackc/pgx/v5"
)

// Obtains a new connection to the Postgres database
func NewPostgresDatabase(config *ServerBackendConfig) (*postGres.Conn, error) {
	if config == nil {
		return nil, fmt.Errorf("No Postgres configuration provided.")
	}

	if config.Url != "" {
		var connectionUrl string = "postgres://%s:%s@%s:%d/%s"

		connectionUrl = fmt.Sprintf(connectionUrl,
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			"zanalytics", // Default database name
		)

		config.Url = connectionUrl
	}

	conn, err := postGres.Connect(context.Background(), config.Url)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return conn, nil
}

// func createTable(name string, conn *postGres.Conn) {}

// func InserValue(name string, conn *postGres.Conn) {}

// func InsertValues(name string, values []any{}, conn *postGres.Conn) {}
