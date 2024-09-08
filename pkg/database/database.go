package database

import (
	"github.com/kicodelibrary/go-http-server-2024/api"
	"github.com/kicodelibrary/go-http-server-2024/pkg/database/errors"
	"github.com/kicodelibrary/go-http-server-2024/pkg/database/mock"
)

// Config holds the database configuration.
type Config struct {
	Type string
}

// NewUsers generates an implementation from the configuration.
func (config Config) NewUsers() (Users, error) {
	switch config.Type {
	case "mock":
		return mock.NewUsers(), nil
	default:
		return nil, errors.ErrInvalidDatabaseType
	}
}

// Users is the interface that wraps the basic user database operations.
type Users interface {
	// List lists all users.
	List() ([]api.User, error)
	// Create creates a new user.
	Create(user api.User) error
	// Get gets a single user with the given ID.
	Get(id string) (api.User, error)
	// Update updates a user.
	Update(id string, user api.User) error
	// Delete deletes a user.
	Delete(id string) error
}
