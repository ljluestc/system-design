package main

import (
    "messaging_system/api"
    "messaging_system/services"
    "messaging_system/utils"
)

func main() {
    log := utils.NewLogger()
    msgService := services.NewMessagingService()
    api.StartServer("8080", msgService, log)
}