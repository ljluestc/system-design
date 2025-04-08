package services

import (
    "database-demo/utils"
)

func GetShardID(userID string) int {
    return utils.HashToInt(userID) % len(dbPools)
}

func InsertToShard(shardID int, userID, data string) error {
    db := dbPools[shardID]
    _, err := db.Exec("INSERT INTO records (user_id, data) VALUES ($1, $2)", userID, data)
    return err
}

func QueryShard(shardID int, userID string) (map[string]string, error) {
    db := dbPools[shardID]
    row := db.QueryRow("SELECT data FROM records WHERE user_id = $1", userID)
    var data string
    err := row.Scan(&data)
    if err != nil {
        return nil, err
    }
    return map[string]string{"user_id": userID, "data": data}, nil
}