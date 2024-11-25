package startup

import "os"

type Config struct {
	Port          string
	ProjectDBHost string
	ProjectDBPort string
	UserHost      string
	UserPort      string
}

func NewConfig() *Config {
	return &Config{
		Port:          os.Getenv("PROJECT_SERVICE_PORT"),
		ProjectDBHost: os.Getenv("PROJECT_DB_HOST"),
		ProjectDBPort: os.Getenv("PROJECT_DB_PORT"),
		UserHost:      os.Getenv("USER_SERVICE_HOST"),
		UserPort:      os.Getenv("USER_SERVICE_PORT"),
	}
}
