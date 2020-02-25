package main

import (
	"log"

	"github.com/VolticFroogo/QShrtn/db"
	"github.com/VolticFroogo/QShrtn/handle"
	"github.com/VolticFroogo/QShrtn/helper"
)

func main() {
	// Seed the random number generator.
	helper.Seed()

	// Initialise the DB.
	err := db.Init()
	if err != nil {
		log.Print(err)
		return
	}

	// Start handling incoming requests.
	err = handle.Start()
	if err != nil {
		log.Print(err)
		return
	}
}
