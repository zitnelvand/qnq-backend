package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultHTTPPort = "3000"

type HTTPConfig struct {
	Host string
	Port string
}

func (c HTTPConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func LoadHTTPConfig() (HTTPConfig, error) {
	if err := loadDotEnv(".env"); err != nil {
		return HTTPConfig{}, err
	}

	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultHTTPPort
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return HTTPConfig{}, fmt.Errorf("HTTP_PORT must be an integer between 1 and 65535, got %q", port)
	}

	return HTTPConfig{Host: host, Port: port}, nil
}

func loadDotEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load dotenv file %q: %w", path, err)
	}

	return nil
}
