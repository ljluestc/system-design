package config

type Config struct {
    NumShards int
}

var globalConfig *Config

func LoadConfig() *Config {
    if globalConfig == nil {
        globalConfig = &Config{NumShards: 4} // Default
    }
    return globalConfig
}

func GetConfig() *Config {
    return LoadConfig()
}