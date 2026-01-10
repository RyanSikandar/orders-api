package application

import (
	"fmt"
	"os"
)

type Config struct {
	RedisAddress string
	ServerPort   uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   3000,
	}

	// Load from environment variables
	if redisAddr, exists := os.LookupEnv("REDIS_ADDRESS"); exists {
		cfg.RedisAddress = redisAddr
	}
	
	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		// Convert serverPort to uint16
		var port uint16
		_, err := fmt.Sscan(serverPort, &port)
		if err == nil {
			cfg.ServerPort = port
		}
	}

	return cfg
}
