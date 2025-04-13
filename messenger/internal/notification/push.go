package notification

import (
    "errors"
    "messenger/internal/utils"
)

type Service struct {
    fcmKey string
}

func NewService(fcmKey string) *Service {
    return &Service{fcmKey: fcmKey}
}

func (s *Service) SendPush(userID, message string) error {
    // Validate input
    if userID == "" || message == "" {
        return errors.New("userID and message required")
    }
    // Pseudocode: Send via Firebase Cloud Messaging
    return utils.SendFCMNotification(s.fcmKey, userID, message)
}