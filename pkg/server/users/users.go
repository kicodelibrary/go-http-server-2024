package users

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kicodelibrary/go-http-server-2024/api"
	"github.com/kicodelibrary/go-http-server-2024/pkg/database"
	dbErrors "github.com/kicodelibrary/go-http-server-2024/pkg/database/errors"
)

// Handler handles the `/users/` routes.
type Handler struct {
	users database.Users
}

// New creates a new handler.
func New(users database.Users) *Handler {
	return &Handler{
		users: users,
	}
}

// AddRoutes adds routes dynamically to the router.
// The argument passed would be a sub-router with the prefix `/users`.
func (h Handler) AddRoutes(r *mux.Router) {
	// List users (GET requests on the `/users/` route.)
	r.HandleFunc("/", h.List).Methods("GET")

	// Create users (POST request to /users/).
	r.HandleFunc("/", h.Create).Methods("POST")

	// Get users (GET request to /users/{id}).
	// {id} is a variable path (not a query).
	r.HandleFunc("/{id}", h.Get).Methods("GET")

	// Update users (PUT request to /users/{id}).
	r.HandleFunc("/{id}", h.Update).Methods("PUT")

	// Delete users (DELETE request to /users/{id}).
	r.HandleFunc("/{id}", h.Delete).Methods("DELETE")
}

// List handles the list user route (`/`).
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	// The response is always going to be JSON.
	w.Header().Set("Content-Type", "application/json")

	// List users from the database.
	users, err := h.users.List()
	if err != nil {
		log.Printf("could not list users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		// Sometimes we may want to return the actual error to the user.
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	msg, err := json.Marshal(users)
	if err != nil {
		log.Printf("could not marshal response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		// Sometimes we may want to return the actual error to the user.
		w.Write(api.NewJSONResponse("internal error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

// Create handles the create user route (`/`).
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	// The response is always going to be JSON.
	w.Header().Set("Content-Type", "application/json")

	// Check the content-type header.
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write(api.NewJSONResponse("Content-Type must be application/json"))
		return
	}

	// Parse and validate the JSON request.
	var user api.User

	// Read the body and decode it.
	body, err := io.ReadAll(r.Body)
	// Always close the body after reading it.
	defer r.Body.Close()
	if err != nil {
		log.Printf("could not decode body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		// Sometimes we may want to return the actual error to the user.
		w.Write(api.NewJSONResponse("Unable to parse the request body"))
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("could not unmarshal JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		// Sometimes we may want to return the actual error to the user.
		w.Write(api.NewJSONResponse("Unable to unmarshal JSON"))
		return
	}

	// Validate the request.
	if err := user.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewJSONResponse("invalid request"))
		return
	}

	// Check if user already exists and error if true.
	_, err = h.users.Get(user.ID)
	if err == nil {
		// This means that the user already exists.
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewJSONResponse("user already exists"))
		return
	} else if !errors.Is(err, dbErrors.ErrUserNotFound) {
		// The error is not errors.ErrUserNotFound. This means that something else went wrong.
		log.Printf("could not get user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	// Create the user.
	err = h.users.Create(user)
	if err != nil {
		log.Printf("could not create user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("user could not be created"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(api.NewJSONResponse("user created"))
}

// Get gets a user. It handles a GET request for the dynamic route `/users/{id}`.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// The response is always going to be JSON.
	w.Header().Set("Content-Type", "application/json")

	// Get the ID from the path.
	id, ok := mux.Vars(r)["id"] // Don't use brackets here (`{}`).
	if !ok {
		// This is mostly a problem with the code.
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	// Check if the user exists.
	user, err := h.users.Get(id)
	if err != nil {
		if errors.Is(err, dbErrors.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(api.NewJSONResponse("user not found"))
			return
		}
		log.Printf("could not get user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}
	// Marshal the user as a JSON and return to the client.
	msg, err := json.Marshal(user)
	if err != nil {
		log.Printf("could not unmarshal the user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

// Update updates the user (PUT request).
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// The response is always going to be JSON.
	w.Header().Set("Content-Type", "application/json")

	// Check the content-type header.
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write(api.NewJSONResponse("Content-Type must be application/json"))
		return
	}

	// Get the ID from the path.
	id, ok := mux.Vars(r)["id"] // Don't use brackets here (`{}`).
	if !ok {
		// This is mostly a problem with the code.
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	// Check and error if user is not found.
	// If err == nil, the user exists and we proceed.
	// If not, we check the error type and return an appropriate code.
	_, err := h.users.Get(id)
	if err != nil {
		if errors.Is(err, dbErrors.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(api.NewJSONResponse("user not found"))
			return
		}
		log.Printf("could not get user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	// Parse the request body and get the update.
	var update api.User

	// Read the body and decode it.
	body, err := io.ReadAll(r.Body)
	// Always close the body after reading it.
	defer r.Body.Close()
	if err != nil {
		log.Printf("could not decode body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		// Sometimes we may want to return the actual error to the user.
		w.Write(api.NewJSONResponse("Unable to parse the request body"))
		return
	}
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("could not unmarshal JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		// Sometimes we may want to return the actual error to the user.
		w.Write(api.NewJSONResponse("Unable to unmarshal JSON"))
		return
	}

	// Validate the request. Check that the user ID is correct.
	if err := update.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewJSONResponse("invalid request"))
		return
	}
	if id != update.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewJSONResponse("ID in the body does not match the path"))
		return
	}

	if err := h.users.Update(id, update); err != nil {
		log.Printf("could not update user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error, could not update user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(api.NewJSONResponse("user updated"))
}

// Delete deletes the user.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// The response is always going to be JSON.
	w.Header().Set("Content-Type", "application/json")

	// Get the ID from the path.
	id, ok := mux.Vars(r)["id"] // Don't use brackets here (`{}`).
	if !ok {
		// This is mostly a problem with the code.
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	// Check and error if user is not found.
	// If err != nil, check the error type and return an appropriate response.
	// If err == nil, the user exists, so proceed to deletion.
	_, err := h.users.Get(id)
	if err != nil {
		if errors.Is(err, dbErrors.ErrUserNotFound) {
			log.Printf("could not get user: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(api.NewJSONResponse("user not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error"))
		return
	}

	// Delete the user.
	if err := h.users.Delete(id); err != nil {
		log.Printf("could not delete user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewJSONResponse("internal error, could not delete user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(api.NewJSONResponse("user deleted"))
}
