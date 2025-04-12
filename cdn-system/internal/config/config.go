package config

import (
    "errors"
    "fmt"
    "log"
    "os"
    "strconv"
    "sync"
)

// Config defines the CDN system configuration
type Config struct {
    Port          string
    EdgeServers   int
    CacheTTL      int
    OriginHost    string
    MaxConnections int
    LogLevel       string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    cfg := &Config{
        Port:          getEnv("PORT", defaults.DefaultPort),
        EdgeServers:   getEnvAsInt("EDGE_SERVERS", defaults.DefaultEdgeServers),
        CacheTTL:      getEnvAsInt("CACHE_TTL", defaults.DefaultCacheTTL),
        OriginHost:    getEnv("ORIGIN_HOST", defaults.DefaultOriginHost),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", defaults.DefaultMaxConnections),
        LogLevel:       getEnv("LOG_LEVEL", defaults.DefaultLogLevel),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}

// getEnv retrieves an environment variable or returns a default
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

// getEnvAsInt converts an environment variable to an integer
func getEnvAsInt(key string, defaultValue int) int {
    if valueStr, exists := os.LookupEnv(key); exists {
        if value, err := strconv.Atoi(valueStr); err == nil {
            return value
        }
    }
    return defaultValue
}

// validateConfig ensures configuration values are valid
func validateConfig(cfg *Config) error {
    if cfg.Port == "" {
        return errors.New("port cannot be empty")
    }
    if cfg.EdgeServers <= 0 {
        return errors.New("edge servers must be positive")
    }
    if cfg.CacheTTL < 0 {
        return errors.New("cache TTL cannot be negative")
    }
    if cfg.OriginHost == "" {
        return errors.New("origin host cannot be empty")
    }
    if cfg.MaxConnections < 1 {
        return errors.New("max connections must be positive")
    }
    if cfg.LogLevel != "info" && cfg.LogLevel != "debug" && cfg.LogLevel != "error" {
        return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
    }
    return nil
}

// Singleton pattern implementation
var (
    once     sync.Once
    instance *Config
)

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
    once.Do(func() {
        var err error
        instance, err = LoadConfig()
        if err != nil {
            log.Fatalf("Failed to load configuration: %v", err)
        }
    })
    return instance
}

// Expand with additional configuration management
// ...