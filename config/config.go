package config

import (
	"encoding/json"
	"log"
	"os"
)

// Load loads a generic config file.
func Load(location string, cfg interface{}) (err error) {
	// Log that we are loading the config file.
	log.Printf("Loading config file: %v", location)

	// Open the config file.
	file, err := os.Open(location)
	if err != nil {
		return
	}

	// Parse the config file's JSON into our generic interface.
	err = json.NewDecoder(file).Decode(&cfg)
	return
}
