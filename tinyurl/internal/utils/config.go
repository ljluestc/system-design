package utils

import "os"

// Config holds system configuration
type Config struct {
    Port          string
    MongoURI      string
    MemcachedHost string
}

// NewConfig loads configuration
func NewConfig() *Config {
    return &Config{
        Port:          getEnv("PORT", "8080"),
        MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
        MemcachedHost: getEnv("MEMCACHED_HOST", "localhost:11211"),
    }
}

// getEnv retrieves an environment variable or fallback
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}