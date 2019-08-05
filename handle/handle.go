package handle

import (
	"log"
	"net/http"

	"github.com/VolticFroogo/QShrtn/config"
	"github.com/VolticFroogo/QShrtn/redirect"
	"github.com/gorilla/mux"
)

const (
	configDirectory = "configs/handle.ini"
)

// Config is the config structure.
type Config struct {
	Port, Certificate, Key string
	SSL                    bool
}

// Start begins listening for all incoming requests.
func Start() {
	// Load the config.
	cfg := Config{}
	err := config.Load(configDirectory, &cfg)
	if err != nil {
		log.Print(err)
		return
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Handle("/new/", http.HandlerFunc(redirect.New)).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./static/"))

	// Handle all static files.
	r.Path("/").Handler(fileServer)
	r.Path("/robots.txt").Handler(fileServer)
	r.PathPrefix("/not-found/").Handler(fileServer)
	r.PathPrefix("/img/").Handler(fileServer)
	r.PathPrefix("/css/").Handler(fileServer)
	r.PathPrefix("/js/").Handler(fileServer)

	r.Handle("/{id}", http.HandlerFunc(redirect.Handle))

	if cfg.SSL {
		// If we are using SSL encryption (HTTPS):
		log.Printf("Listening for incoming HTTPS requests on port %v.", cfg.Port)

		// Serve TLS using the certificate and key files from the config.
		http.ListenAndServeTLS(":"+cfg.Port, cfg.Certificate, cfg.Key, r)
	} else {
		// Otherwise:
		log.Printf("Listening for incoming HTTP requests on port %v.", cfg.Port)

		// Serve plain HTTP responses.
		http.ListenAndServe(":"+cfg.Port, r)
	}
}
