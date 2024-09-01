package users

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kicodelibrary/go-http-server-2024/api"
)

// Handler handles the `/users/` routes.
type Handler struct {
	users map[string]api.User // In-memory storage. Will be replaced when there's a database.
}

// New creates a new handler.
func New() *Handler {
	return &Handler{
		users: make(map[string]api.User),
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

	var users []api.User
	for _, v := range h.users {
		users = append(users, v)
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
	_, ok := h.users[user.ID]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewJSONResponse("user already exists"))
		return
	}
	// Create the user (add it to the map).
	h.users[user.ID] = user
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
	user, ok := h.users[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(api.NewJSONResponse("user not found"))
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

	// Check if the user exists.
	_, ok = h.users[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(api.NewJSONResponse("user not found"))
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

	h.users[id] = update
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

	// Check if the user exists.
	_, ok = h.users[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(api.NewJSONResponse("user not found"))
		return
	}

	// Delete the user.
	delete(h.users, id)
	w.WriteHeader(http.StatusOK)
	w.Write(api.NewJSONResponse("user deleted"))
}
