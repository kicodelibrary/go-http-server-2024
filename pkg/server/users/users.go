package users

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler handles the `/users/` routes.
type Handler struct {
}

// AddRoutes adds routes dynamically to the router.
// The argument passed would be a sub-router with the prefix `/users`.
func (h Handler) AddRoutes(r *mux.Router) {
	// Handle GET requests on the `/users/` route.
	r.HandleFunc("/", h.List).Methods("GET")

	// Add routes for Create, Read, Update, Delete here.
}

// List handles the list user route (`/`).
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "List Users")
}
