package backend

// ServerBackendConfig represents the configuration for a specific backend in the YAML file
type ServerBackendConfig struct {
	Url string "json:\"url\" yaml:\"url\""
}

// ServerBackendsConfig represents the backend configuration in the YAML file
type ServerBackendsConfig struct {
	Database struct {
		Client string "json:\"client\" yaml:\"client\""
	}
	Postgres *ServerBackendConfig `json:"postgres" yaml:"postgres"`
	Redis    *ServerBackendConfig `json:"redis" yaml:"redis"`
	RabbitMQ *ServerBackendConfig `json:"rabbitmq" yaml:"rabbitmq"`
}

// YamlConfig represents the structure of the configuration YAML file
type ServerConfig struct {
	Config struct {
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
		Domains  []string              "json:\"domains\" yaml:\"domains\""
		Backends *ServerBackendsConfig `json:"backends" yaml:"backends"`
	}
}
