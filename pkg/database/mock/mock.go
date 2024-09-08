// Package mock provides an in-memory mock implementation of the database interfaces.
package mock

import (
	"github.com/kicodelibrary/go-http-server-2024/api"
	"github.com/kicodelibrary/go-http-server-2024/pkg/database/errors"
)

// Users mocks users.
type Users struct {
	users map[string]api.User
}

// NewUsers returns a new mock users.
func NewUsers() *Users {
	return &Users{
		users: make(map[string]api.User),
	}
}

// List implements database.Users.
func (u Users) List() ([]api.User, error) {
	ret := []api.User{} // Initialize.
	for _, user := range u.users {
		ret = append(ret, user)
	}
	return ret, nil
}

// Create implements database.Users.
func (u Users) Create(user api.User) error {
	u.users[user.ID] = user
	return nil
}

// Get implements database.Users.
// If user exists, there is no error.
// If user does not exists this function returns database.ErrUserNotFound.
func (u Users) Get(id string) (api.User, error) {
	user, ok := u.users[id]
	if !ok {
		return api.User{}, errors.ErrUserNotFound
	}
	return user, nil
}

// Update implements database.Users.
func (u Users) Update(id string, user api.User) error {
	u.users[id] = user // This is a replacement.
	return nil
}

// Delete implements database.Users.
func (u Users) Delete(id string) error {
	delete(u.users, id)
	return nil
}
