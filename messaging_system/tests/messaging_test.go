package tests

import (
    "testing"
    "messaging_system/services"
)

func TestSendMessage(t *testing.T) {
    msgService := services.NewMessagingService()
    err := msgService.SendMessage("Hello, World!")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}