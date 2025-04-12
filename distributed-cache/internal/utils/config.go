package utils

type Config struct {
    ServerPort string
    Nodes      []string
}

func LoadConfig() *Config {
    // Hardcoded for simplicity; replace with file/env loading in production
    return &Config{
        ServerPort: "8080",
        Nodes:      []string{"node1", "node2", "node3"},
    }
}