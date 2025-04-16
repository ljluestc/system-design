package config

var defaultConfig = &Config{NumEdges: 4, CacheCapacity: 100}

func GetConfig() *Config {
    return defaultConfig
}