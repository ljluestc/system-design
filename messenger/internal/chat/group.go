package chat

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "messenger/pkg/db"
)

type Group struct {
    ID        string
    Name      string
    Members   []string
    CreatedAt time.Time
}

func (s *Service) CreateGroup(name, creatorID string) (Group, error) {
    // Validate input
    if name == "" || creatorID == "" {
        return Group{}, errors.New("name and creatorID required")
    }
    group := Group{
        ID:        uuid.New().String(),
        Name:      name,
        Members:   []string{creatorID},
        CreatedAt: time.Now(),
    }
    // Save to database
    if err := s.db.SaveGroup(group); err != nil {
        return Group{}, err
    }
    return group, nil
}

func (s *Service) AddMember(groupID, userID string) error {
    return s.db.AddGroupMember(groupID, userID)
}