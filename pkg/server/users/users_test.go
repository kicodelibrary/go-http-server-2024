package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kicodelibrary/go-http-server-2024/api"
	. "github.com/kicodelibrary/go-http-server-2024/pkg/server/users"
	"github.com/smarty/assertions"
)

func TestUser(t *testing.T) {
	a := assertions.New(t)
	h := New()

	// Create a test router.
	router := mux.NewRouter().PathPrefix("/users").Subrouter()
	h.AddRoutes(router)

	// Create the request.
	req, err := http.NewRequest(http.MethodGet, "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test HTTP recorder.
	// This prevents us from needing to create an HTTP Server for testing.
	// The recorder simulates a server and captures the response.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response := rec.Result()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(string(body), assertions.ShouldEqual, "[]")

	// Test the create user function.
	alice := api.User{
		ID:   "alice",
		Name: "Alice",
		Age:  30,
	}
	aliceMsg, err := json.Marshal(alice)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(aliceMsg))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response = rec.Result()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(response.StatusCode, assertions.ShouldEqual, http.StatusCreated)
	a.So(string(body), assertions.ShouldEqual, `{"message":"user created"}`)

	// Test the get user method.
	req, err = http.NewRequest(http.MethodGet, "/users/alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response = rec.Result()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(response.StatusCode, assertions.ShouldEqual, http.StatusOK)
	a.So(body, assertions.ShouldEqual, aliceMsg)

	// Test the update method.
	alice.Name = "Alice Smith"
	aliceMsgUpdated, err := json.Marshal(alice)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(http.MethodPut, "/users/alice", bytes.NewBuffer(aliceMsgUpdated))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response = rec.Result()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(response.StatusCode, assertions.ShouldEqual, http.StatusOK)
	a.So(string(body), assertions.ShouldEqual, `{"message":"user updated"}`)

	// Get the updated user.
	req, err = http.NewRequest(http.MethodGet, "/users/alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response = rec.Result()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(response.StatusCode, assertions.ShouldEqual, http.StatusOK)
	a.So(body, assertions.ShouldEqual, aliceMsgUpdated)

	// Test the delete method.
	req, err = http.NewRequest(http.MethodDelete, "/users/alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response = rec.Result()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(response.StatusCode, assertions.ShouldEqual, http.StatusOK)
	a.So(string(body), assertions.ShouldEqual, `{"message":"user deleted"}`)

	// Get after deletion.
	req, err = http.NewRequest(http.MethodGet, "/users/alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	response = rec.Result()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		response.Body.Close()
	})
	a.So(response.StatusCode, assertions.ShouldEqual, http.StatusNotFound)
	a.So(string(body), assertions.ShouldEqual, `{"message":"user not found"}`)
}

func TestUsers(t *testing.T) {
	a := assertions.New(t)
	h := New()

	// Create a test router.
	router := mux.NewRouter().PathPrefix("/users").Subrouter()
	h.AddRoutes(router)

	// Use a range loop to test the methods.
	// Define an anonymous struct for the test parameters.
	for _, tc := range []struct {
		// Definition
		Name         string
		Request      func() *http.Request // This gives us more flexibility to define things.
		ResponseCode int
		ResponseBody string
	}{
		// Test cases
		{
			Name: "List",
			Request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/users/", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			ResponseCode: http.StatusOK,
			ResponseBody: "[]",
		},
		{
			Name: "CreateIncorrectContentType",
			Request: func() *http.Request {
				alice := api.User{
					ID:   "alice",
					Name: "Alice",
					Age:  30,
				}
				aliceMsg, err := json.Marshal(alice)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(aliceMsg))
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			ResponseCode: http.StatusUnsupportedMediaType,
			ResponseBody: `{"message":"Content-Type must be application/json"}`,
		},
		{
			Name: "Create",
			Request: func() *http.Request {
				alice := api.User{
					ID:   "alice",
					Name: "Alice",
					Age:  30,
				}
				aliceMsg, err := json.Marshal(alice)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(aliceMsg))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			ResponseCode: http.StatusCreated,
			ResponseBody: `{"message":"user created"}`,
		},
		{
			Name: "Get",
			Request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/users/alice", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			ResponseCode: http.StatusOK,
			ResponseBody: `{"id":"alice","name":"Alice","age":30}`,
		},
		{
			Name: "Update",
			Request: func() *http.Request {
				alice := api.User{
					ID:   "alice",
					Name: "Alice Smith",
					Age:  30,
				}
				aliceMsgUpdated, err := json.Marshal(alice)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(http.MethodPut, "/users/alice", bytes.NewBuffer(aliceMsgUpdated))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			ResponseCode: http.StatusOK,
			ResponseBody: `{"message":"user updated"}`,
		},
		{
			Name: "GetAfterUpdate",
			Request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/users/alice", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			ResponseCode: http.StatusOK,
			ResponseBody: `{"id":"alice","name":"Alice Smith","age":30}`,
		},
		{
			Name: "Delete",
			Request: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, "/users/alice", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			ResponseCode: http.StatusOK,
			ResponseBody: `{"message":"user deleted"}`,
		},
		{
			Name: "GetAfterDelete",
			Request: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, "/users/alice", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			ResponseCode: http.StatusNotFound,
			ResponseBody: `{"message":"user not found"}`,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, tc.Request())
			response := rec.Result()
			if !a.So(response.StatusCode, assertions.ShouldEqual, tc.ResponseCode) {
				t.Fatalf("unexpected status code: %d", response.StatusCode)
			}
			// Read the response body.
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				response.Body.Close()
			})
			if !a.So(string(body), assertions.ShouldEqual, tc.ResponseBody) {
				t.Fatalf("unexpected response body: %s", string(body))
			}
		})
	}

}

func TestUsersMulti(t *testing.T) {
	a := assertions.New(t)
	h := New()

	// Create a test router.
	router := mux.NewRouter().PathPrefix("/users").Subrouter()
	h.AddRoutes(router)

	// Create an outer loop that feeds the inner loop with users.
	for _, tu := range []struct {
		UserName string
		User     api.User
	}{
		{
			UserName: "alice",
			User: api.User{
				ID:   "alice",
				Age:  30,
				Name: "Alice",
			},
		},
		{
			UserName: "bob",
			User: api.User{
				ID:   "bob",
				Age:  25,
				Name: "Bob",
			},
		},
	} {
		// Use a range loop to test the methods.
		// Define an anonymous struct for the test parameters.
		for _, tc := range []struct {
			// Definition
			Name             string
			RequestFunc      func() *http.Request // This gives us more flexibility to define things.
			ResponseCode     int
			ResponseBodyFunc func() string
		}{
			// Test cases
			{
				Name: "List",
				RequestFunc: func() *http.Request {
					req, err := http.NewRequest(http.MethodGet, "/users/", nil)
					if err != nil {
						t.Fatal(err)
					}
					return req
				},
				ResponseCode:     http.StatusOK,
				ResponseBodyFunc: func() string { return "[]" },
			},
			{
				Name: "CreateIncorrectContentType",
				RequestFunc: func() *http.Request {
					userMsg, err := json.Marshal(tu.User)
					if err != nil {
						t.Fatal(err)
					}
					req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(userMsg))
					if err != nil {
						t.Fatal(err)
					}
					return req
				},
				ResponseCode: http.StatusUnsupportedMediaType,
				ResponseBodyFunc: func() string {
					return `{"message":"Content-Type must be application/json"}`
				},
			},
			{
				Name: "Create",
				RequestFunc: func() *http.Request {
					userMsg, err := json.Marshal(tu.User)
					if err != nil {
						t.Fatal(err)
					}
					req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(userMsg))
					if err != nil {
						t.Fatal(err)
					}
					req.Header.Set("Content-Type", "application/json")
					return req
				},
				ResponseCode: http.StatusCreated,
				ResponseBodyFunc: func() string {
					return `{"message":"user created"}`
				},
			},
			{
				Name: "Get",
				RequestFunc: func() *http.Request {
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", tu.UserName), nil)
					if err != nil {
						t.Fatal(err)
					}
					return req
				},
				ResponseCode: http.StatusOK,
				ResponseBodyFunc: func() string {
					msg, err := json.Marshal(tu.User)
					if err != nil {
						t.Fatal(err)
					}
					return string(msg)
				},
			},
			{
				Name: "Update",
				RequestFunc: func() *http.Request {
					userUpdated := tu.User
					userUpdated.Name = "Abcd"
					userMsgUpdated, err := json.Marshal(userUpdated)
					if err != nil {
						t.Fatal(err)
					}
					req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", tu.UserName), bytes.NewBuffer(userMsgUpdated))
					if err != nil {
						t.Fatal(err)
					}
					req.Header.Set("Content-Type", "application/json")
					return req
				},
				ResponseCode: http.StatusOK,
				ResponseBodyFunc: func() string {
					return `{"message":"user updated"}`
				},
			},
			{
				Name: "GetAfterUpdate",
				RequestFunc: func() *http.Request {
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", tu.UserName), nil)
					if err != nil {
						t.Fatal(err)
					}
					return req
				},
				ResponseCode: http.StatusOK,
				ResponseBodyFunc: func() string {
					updated := tu.User
					updated.Name = "Abcd"
					msg, err := json.Marshal(updated)
					if err != nil {
						t.Fatal(err)
					}
					return string(msg)
				},
			},
			{
				Name: "Delete",
				RequestFunc: func() *http.Request {
					req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", tu.UserName), nil)
					if err != nil {
						t.Fatal(err)
					}
					return req
				},
				ResponseCode: http.StatusOK,
				ResponseBodyFunc: func() string {
					return `{"message":"user deleted"}`
				},
			},
			{
				Name: "GetAfterDelete",
				RequestFunc: func() *http.Request {
					req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", tu.UserName), nil)
					if err != nil {
						t.Fatal(err)
					}
					return req
				},
				ResponseCode: http.StatusNotFound,
				ResponseBodyFunc: func() string {
					return `{"message":"user not found"}`
				},
			},
		} {
			t.Run(fmt.Sprintf("%s/%s", tu.UserName, tc.Name), func(t *testing.T) {
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, tc.RequestFunc())
				response := rec.Result()
				if !a.So(response.StatusCode, assertions.ShouldEqual, tc.ResponseCode) {
					t.Fatalf("unexpected status code: %d", response.StatusCode)
				}
				// Read the response body.
				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					response.Body.Close()
				})
				if !a.So(string(body), assertions.ShouldEqual, tc.ResponseBodyFunc()) {
					t.Fatalf("unexpected response body: %s", string(body))
				}
			})
		}
	}
}
