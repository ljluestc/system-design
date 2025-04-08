package services

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

var dbPools = make(map[int]*sql.DB)

func InitDatabase() {
    for i := 0; i < 3; i++ {
        connStr := fmt.Sprintf("user=postgres password=password dbname=demo host=shard-%d port=5432 sslmode=disable", i)
        db, err := sql.Open("postgres", connStr)
        if err != nil {
            log.Fatal(err)
        }
        dbPools[i] = db
    }
}