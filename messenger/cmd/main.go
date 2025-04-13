package main

import (
    "github.com/calelin/messenger/api"
    "github.com/calelin/messenger/config"
    "github.com/calelin/messenger/internal/auth"
    "github.com/calelin/messenger/internal/chat"
    "github.com/calelin/messenger/internal/notification"
    "github.com/calelin/messenger/internal/utils"
    "github.com/calelin/messenger/pkg/db"
)

func main() {
    log := utils.NewLogger()
    cfg := config.NewConfig()
    dbConn := db.NewCassandraConn(cfg)
    authService := auth.NewAuthService(dbConn, log)
    chatService := chat.NewChatService(dbConn, log)
    notificationService := notification.NewPushService(log)

    // Start the server
    api.StartServer(cfg, authService, chatService, notificationService, log)
}