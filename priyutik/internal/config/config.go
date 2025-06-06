package config

type Config struct {
	DB struct {
		URI string // Добавляем поле для URI подключения
	}
	Server struct {
		Port string
	}
	Session struct {
		SecretKey string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Database configuration
	cfg.DB.URI = "postgresql://postgres:0000@localhost:5432/priyutik?sslmode=disable"

	// Server configuration
	cfg.Server.Port = "8080"

	// Session configuration
	cfg.Session.SecretKey = "secret-key"

	return cfg, nil
}
