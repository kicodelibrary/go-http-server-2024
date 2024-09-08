package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kicodelibrary/go-http-server-2024/pkg/database"
	"github.com/kicodelibrary/go-http-server-2024/pkg/server/users"
	"github.com/spf13/pflag"
)

// Config is the configuration for the application.
type Config struct {
	Port, Host string
	Timeout    time.Duration
	Database   database.Config
}

var (
	config = &Config{} // This holds the configuration.
	flags  = pflag.NewFlagSet("server", pflag.ExitOnError)
)

func main() {
	// Define the root router.
	root := mux.NewRouter()

	// Handle the default home (index) route.
	// This only works for GET.
	// For other methods from the client on this route, the server will return an error.
	root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	}).Methods("GET")

	// Handle the `/users` routes.
	usersDB, err := config.Database.NewUsers()
	if err != nil {
		log.Fatal(err)
	}
	h := users.New(usersDB)

	// Create a subrouter for the `/users` prefix.
	sub := root.PathPrefix("/users").Subrouter()
	h.AddRoutes(sub)

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	server := &http.Server{
		Addr:           address,
		Handler:        root,
		ReadTimeout:    config.Timeout,
		WriteTimeout:   config.Timeout,
		MaxHeaderBytes: 1 << 20, // Restrict the max size of headers.
	}

	log.Printf("Start server: %s\n", address)

	// Run the server.
	log.Fatal(server.ListenAndServe())
}

// init gets called before main().
func init() {
	// Define the flags for the server.
	flags.StringVarP(&config.Host, "host", "h", "localhost", "Hostname")
	flags.StringVarP(&config.Port, "port", "p", "8080", "Port")
	flags.DurationVarP(&config.Timeout, "timeout", "t", 10*time.Second, "Server timeouts")

	// Define the flags for the database.
	flags.StringVar(&config.Database.Type, "database.type", "mock", "Database type (supported values: mock)")

	// Define the usage (help) function (when `--help` is used).
	flags.Usage = func() {
		usage := `Usage: server [flags]
An HTTP Server to manage users.
`
		// Print this message at the top.
		fmt.Fprintln(os.Stderr, usage)
		fmt.Fprintln(os.Stderr, flags.FlagUsages())
	}

	// Parse the flags.
	// This is important to actually read the values into config.
	// Args[0] is always the name of the command.
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("could not parse flags: %v", err)
	}

}
