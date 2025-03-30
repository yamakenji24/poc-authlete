package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthleteBaseURL      string
	AuthleteServiceID    string
	AuthleteClientID     string
	AuthleteClientSecret string
	AuthleteRedirectURI  string
	AuthleteAccessToken  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		AuthleteBaseURL:      os.Getenv("AUTHLETE_BASE_URL"),
		AuthleteServiceID:    os.Getenv("AUTHLETE_SERVICE_ID"),
		AuthleteClientID:     os.Getenv("AUTHLETE_CLIENT_ID"),
		AuthleteClientSecret: os.Getenv("AUTHLETE_CLIENT_SECRET"),
		AuthleteRedirectURI:  os.Getenv("AUTHLETE_REDIRECT_URI"),
		AuthleteAccessToken:  os.Getenv("AUTHLETE_ACCESS_TOKEN"),
	}, nil
}
