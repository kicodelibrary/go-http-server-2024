// Package api provides API messages.
// All messages are interpreted as JSON.
package api

import (
	"fmt"
	"regexp"
)

// User is a user.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	// Extend this message as required.
}

// Validate users.

// ID can only contain lowercase letters and numbers.
// It should be at least 3 characters and max 32.
var userIDRegexp = regexp.MustCompile(`^[a-z0-9]{3,32}$`)

// Validate validates a user based on the conditions.
func (u *User) Validate() error {
	if !userIDRegexp.MatchString(u.ID) {
		return fmt.Errorf("invalid user ID: %s", u.ID)
	}

	// Ex: Add a validation for ages (ex: 18 - 100).
	return nil
}
