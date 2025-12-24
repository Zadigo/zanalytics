package backend

import (
	"context"
	"log"

	postGres "github.com/jackc/pgx/v5"
)

// Obtains a new connection to the Postgres database
func NewPostgresDatabase() *postGres.Conn {
	conn, err := postGres.Connect(context.Background(), "postgres://username:password@localhost:5432/database_name")

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil
	}

	return &conn
}
