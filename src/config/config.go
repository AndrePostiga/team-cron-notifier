package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type HttpClientConfig struct {
	TimeoutInSeconds int `env:"TIMEOUT_IN_SECONDS" envDefault:"10"`
}

type SlackApiConfig struct {
	BaseUrl    string           `env:"SLACK_API_BASE_URL"`
	HTTPClient HttpClientConfig `envPrefix:"SLACK_API_HTTP_CLIENT_"` // SLACK_API_HTTP_CLIENT_TIMEOUT_IN_SECONDS
}

type GithubApiConfig struct {
	BaseUrl    string           `env:"GITHUB_API_BASE_URL"`
	HTTPClient HttpClientConfig `envPrefix:"GITHUB_API_HTTP_CLIENT_"` // GITHUB_API_HTTP_CLIENT_TIMEOUT_IN_SECONDS
}

type Config struct {
	LogLevel   string `env:"LOG_LEVEL"`
	GithubApi  GithubApiConfig
	SlackApi   SlackApiConfig
	HttpClient HttpClientConfig
}

func LoadEnv() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println(".env file does not exist, skipping loading environment variables")
		return nil
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return nil
}

func LoadConfig() (*Config, error) {
	err := LoadEnv()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadConfigGeneric[T any]() (*T, error) {
	err := LoadEnv()
	if err != nil {
		return nil, err
	}

	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
