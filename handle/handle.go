package handle

import (
	"log"
	"net/http"
	"os"

	"github.com/VolticFroogo/QShrtn/redirect"
	"github.com/gorilla/mux"
)

const (
	defaultPort = "80"
)

// Start begins listening for all incoming requests.
func Start() (err error) {
	// Create a new Mux Router with strict slash.
	r := mux.NewRouter()
	r.StrictSlash(true)

	// Handle new URL requests.
	r.Handle("/new/", http.HandlerFunc(redirect.New)).Methods(http.MethodPost)

	// Create a new static file server.
	fileServer := http.FileServer(http.Dir("./static/"))

	// Handle all static files with the file server.
	r.Path("/").Handler(fileServer)
	r.Path("/robots.txt").Handler(fileServer)
	r.Path("/sitemap.xml").Handler(fileServer)
	r.PathPrefix("/not-found/").Handler(fileServer)
	r.PathPrefix("/img/").Handler(fileServer)
	r.PathPrefix("/css/").Handler(fileServer)
	r.PathPrefix("/js/").Handler(fileServer)

	// Handle all unknown links, possibly redirecting links.
	r.Handle("/{id}", http.HandlerFunc(redirect.Handle))
	r.Handle("/{id}/json", http.HandlerFunc(redirect.JSON))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Listening for incoming HTTP requests on port %s.", port)

	// Serve plain HTTP responses.
	err = http.ListenAndServe(":"+port, r)
	return
}
