package config

import (
    "encoding/json"
    "os"
)

// ServiceInstances holds the list of instances for each service
type ServiceInstances struct {
    User         []string `json:"user"`
    Tweet        []string `json:"tweet"`
    Timeline     []string `json:"timeline"`
    Search       []string `json:"search"`
    Notification []string `json:"notification"`
}

// LoadConfig reads the configuration from config.json
func LoadConfig() (*ServiceInstances, error) {
    file, err := os.Open("config.json")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var cfg ServiceInstances
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}