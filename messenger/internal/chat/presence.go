package chat

import (
    "time"
    "messenger/pkg/db"
)

type Presence struct {
    UserID    string
    IsOnline  bool
    LastSeen  time.Time
}

func (s *Service) UpdatePresence(userID string, isOnline bool) error {
    // Update user presence in database
    presence := Presence{
        UserID:   userID,
        IsOnline: isOnline,
        LastSeen: time.Now(),
    }
    return s.db.SavePresence(presence)
}

func (s *Service) GetPresence(userID string) (Presence, error) {
    return s.db.GetPresence(userID)
}