package api

import "testing"

func TestUsers(t *testing.T) {
	// Test validation using range testing.
	for _, test := range []struct {
		// Define the test structure.
		Name           string
		ID             string
		ErrorAssertion func(error) bool
	}{
		// Define the values for test cases.
		{
			Name: "Valid",
			ID:   "abc123",
			ErrorAssertion: func(err error) bool {
				// Since this is a valid ID, there should be no error.
				return err == nil
			},
		},
		{
			Name: "Invalid Length",
			ID:   "ab",
			ErrorAssertion: func(err error) bool {
				// We expect an error here.
				return err != nil
			},
		},
		{
			Name: "Invalid character",
			ID:   "abc123&",
			ErrorAssertion: func(err error) bool {
				// We expect an error here.
				return err != nil
			},
		},
	} {
		// Run the tests.
		t.Run(test.Name, func(t *testing.T) {
			user := &User{
				ID:   test.ID,
				Name: "abcd",
				Age:  20,
			}
			if !test.ErrorAssertion(user.Validate()) {
				t.Errorf("error assertion failed")
			}
		})
	}
}
