package config

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
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
	cfg.DB.Host = "localhost"
	cfg.DB.Port = "5432"
	cfg.DB.User = "postgres"
	cfg.DB.Password = "0000"
	cfg.DB.Name = "priyutik"
	cfg.DB.SSLMode = "disable"

	// Server configuration
	cfg.Server.Port = "8080"

	// Session configuration
	cfg.Session.SecretKey = "secret-key"

	return cfg, nil
}
