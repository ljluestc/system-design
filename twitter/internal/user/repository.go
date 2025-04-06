package user

import (
    "errors"
    "twitter/pkg/database"
)

// UserRepository handles data access for users
type UserRepository struct {
    db *database.MockDB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.MockDB) *UserRepository {
    return &UserRepository{db: db}
}

// Save persists a user to the database
func (repo *UserRepository) Save(user *User) error {
    return repo.db.SaveUser(user)
}

// FindByID retrieves a user by ID
func (repo *UserRepository) FindByID(id string) (*User, error) {
    return repo.db.GetUser(id)
}

// FindByEmail retrieves a user by email
func (repo *UserRepository) FindByEmail(email string) (*User, error) {
    for _, user := range repo.db.Users() { // Access underlying map indirectly
        if user.Email == email {
            return user, nil
        }
    }
    return nil, errors.New("user not found")
}