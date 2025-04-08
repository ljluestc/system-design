package utils

type Config struct {
    Port string
}

func LoadConfig() Config {
    return Config{Port: "3001"}
}