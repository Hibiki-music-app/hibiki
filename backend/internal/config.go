package internal

import (
	"os"
	"strconv"
)

type Config struct {
	Proto   string
	Host    string
	PortInt int
	Port    string
	DSN     string
	User    string
	Passwd  string
}

func Load() *Config {
	return &Config{
		Proto:   getEnv("PROTO", "http"),
		Host:    getEnv("HOST", "localhost"),
		Port:    getEnv("PORT", ":8080"),
		PortInt: getEnvInt("PORTINT", 8080),
		User:    getEnv("USER", "example_user"),
		Passwd:  getEnv("PASSWD", "example_passwd"),
		// TODO: serverTLS config then activate secure https for cookies
		// serverCRT := getEnv("CRT", "")
		// serverKEY := getEnv("KEY", "")
	}
}

// getEnv is a helper function to get the environment variable
func getEnv(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}

func getEnvInt(key string, def int) int {
	if value, exists := os.LookupEnv(key); exists {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return def
}
