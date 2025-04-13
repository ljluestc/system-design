package chat

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "messenger/pkg/db"
)

type Service struct {
    db *db.Cassandra
}

func NewService(db *db.Cassandra) *Service {
    return &Service{db: db}
}

type Message struct {
    ID          string
    SenderID    string
    RecipientID string
    GroupID     string
    Content     string
    Timestamp   time.Time
}

func (s *Service) SendMessage(senderID, recipientID, groupID, content string) (Message, error) {
    // Validate message parameters
    if senderID == "" || (recipientID == "" && groupID == "") {
        return Message{}, errors.New("invalid message parameters")
    }
    msg := Message{
        ID:          uuid.New().String(),
        SenderID:    senderID,
        RecipientID: recipientID,
        GroupID:     groupID,
        Content:     content,
        Timestamp:   time.Now(),
    }
    // Save to database
    if err := s.db.SaveMessage(msg); err != nil {
        return Message{}, err
    }
    return msg, nil
}

func (s *Service) GetMessages(userID string, limit int) ([]Message, error) {
    return s.db.GetMessages(userID, limit)
}