package services

func GetMonitoringStats() map[string]int {
    return map[string]int{"shards": len(dbPools)}
}

func IsHealthy() bool {
    for _, db := range dbPools {
        if err := db.Ping(); err != nil {
            return false
        }
    }
    return true
}