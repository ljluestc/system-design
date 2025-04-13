package auth

// CreateUser creates a new user
func (s *AuthService) CreateUser(username, apiKey string) error {
    if err := s.db.CreateUser(username, apiKey); err != nil {
        s.log.Errorf("Failed to create user %s: %v", username, err)
        return err
    }
    s.log.Infof("User %s created with API key %s", username, apiKey)
    return nil
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(username string) error {
    if err := s.db.DeleteUser(username); err != nil {
        s.log.Errorf("Failed to delete user %s: %v", username, err)
        return err
    }
    s.log.Infof("User %s deleted", username)
    return nil
}