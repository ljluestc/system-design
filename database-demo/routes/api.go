package routes

import (
    "net/http"
    "database-demo/controllers"
)

func SetupAPIRoutes() {
    http.HandleFunc("/api/sharding/insert", controllers.ShardingInsertHandler)
    http.HandleFunc("/api/sharding/get", controllers.ShardingGetHandler)
}