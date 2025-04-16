package config

type Config struct {
    NumEdges, CacheCapacity int
}

func LoadConfig() *Config {
    return &Config{NumEdges: 4, CacheCapacity: 100}
}