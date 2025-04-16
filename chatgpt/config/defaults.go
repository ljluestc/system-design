package config

func SetDefaults() {
    if globalConfig == nil {
        globalConfig = &Config{NumShards: 4}
    }
}