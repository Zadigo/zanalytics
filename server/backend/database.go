package backend

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Zadigo/zanalytics/utils"
	postGres "github.com/jackc/pgx/v5"
)

// Obtains a new connection to the Postgres database
func NewPostgresDatabase(config *ServerBackendConfig) (*postGres.Conn, error) {
	if config == nil {
		return nil, fmt.Errorf("No Postgres configuration provided.")
	}

	if strings.TrimSpace(config.Host) == "" {
		config.Host = "localhost"
		log.Println("No Postgres host provided, defaulting to localhost.")
	}

	if config.Port == 0 || config.Port < 0 {
		config.Port = 5432
		log.Println("No Postgres port provided, defaulting to 5432.")
	}

	// Url matches the format: "postgres://user:password@localhost:5432/zanalytics"
	if strings.TrimSpace(config.Url) == "" {
		var connectionUrl string = "postgres://%s:%s@%s:%d/%s"

		if (config.Username == "") || (config.Password == "") {
			log.Panic("Postgres username or password not provided.")
		}

		connectionUrl = fmt.Sprintf(connectionUrl,
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			"zanalytics", // Default database name
		)

		config.Url = connectionUrl
	}

	// Create a new connection to the Postgres database
	conn, err := postGres.Connect(context.Background(), config.Url)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return conn, nil
}

// Creates necessary tables in the Postgres database
func CreateTables(conn *postGres.Conn, config *ServerBackendConfig) {
	eventsTableQuery := []string{
		"events",
		`
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
	}

	accountsTableQuery := []string{
		"accounts",
		`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
	}

	sqls := [][]string{}
	sqls = append(sqls, eventsTableQuery)
	sqls = append(sqls, accountsTableQuery)

	for _, query := range sqls {
		_, err := conn.Exec(context.Background(), query[1])
		if err != nil {
			fmt.Printf("Error creating %v table: %v\n", query[0], err)
		} else {
			fmt.Printf("%v table created or already exists.\n", query[0])
		}
	}

	err := CreateUser(conn, config.Username, config.Password)

	if err != nil {
		fmt.Printf("Error creating admin user: %v\n", err)
	}
}

// Creates a new user in the Postgres database
func CreateUser(conn *postGres.Conn, username string, password string) error {
	if (username != "") && (password != "") {
		createAdminQuery := `
		INSERT INTO accounts (email, password)
		VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING;
		`

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			fmt.Printf("Error hashing password: %v\n", err)
			return err
		}

		_, err = conn.Exec(context.Background(), createAdminQuery, username, hashedPassword)
		if err != nil {
			fmt.Printf("Error creating admin account: %v\n", err)
		} else {
			fmt.Println("Admin account created or already exists.")
		}
	}

	return nil
}

// A function used to authenticate a user in the Postgres database
func AuthenticateUser(conn *postGres.Conn, username string, password string) bool {
	authQuery := `SELECT password FROM accounts WHERE email=$1;`

	var storedHashedPassword string
	err := conn.QueryRow(context.Background(), authQuery, username).Scan(&storedHashedPassword)

	if err != nil {
		fmt.Printf("Error fetching user data: %v\n", err)
		return false
	}

	isValid := utils.VerifyPassword(password, storedHashedPassword)
	return isValid
}

// func InserValue(name string, conn *postGres.Conn) {}

// func InsertValues(name string, values []any{}, conn *postGres.Conn) {}
