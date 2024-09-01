package users_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kicodelibrary/go-http-server-2024/api"
	. "github.com/kicodelibrary/go-http-server-2024/pkg/server/users"
	"github.com/smarty/assertions"
)

func TestUsers(t *testing.T) {
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
