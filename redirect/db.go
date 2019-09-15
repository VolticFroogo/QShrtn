package redirect

import (
	"database/sql"
	"fmt"

	"github.com/VolticFroogo/QShrtn/helper"

	"github.com/VolticFroogo/QShrtn/db"
	"github.com/VolticFroogo/QShrtn/model"
)

var (
	ErrIDTaken = fmt.Errorf("id already taken")
)

// FromID gets a redirect from the database given an ID.
func FromID(id string) (redirect model.Redirect, err error) {
	// Query a row from our ID.
	row, err := db.Dot.QueryRow(db.SQL, "redirect-from-id", id)
	if err != nil {
		return
	}

	// Scan the row into our profile.
	err = scan(&redirect, row)
	return
}

// FromURL gets a redirect from the database given an URL.
func FromURL(url string) (redirect model.Redirect, err error) {
	// Query a row from our ID.
	row, err := db.Dot.QueryRow(db.SQL, "redirect-from-url", url)
	if err != nil {
		return
	}

	// Scan the row into our profile.
	err = scan(&redirect, row)
	return
}

// Insert a redirect into the database.
func Insert(url string) (redirect model.Redirect, err error) {
	// Check if there's already a redirect for that URL.
	redirect, err = FromURL(url)
	// If there is an entry for this URL already, return that.
	// Or, if there was an error that wasn't no rows, return that.
	if err == nil || err != sql.ErrNoRows {
		return
	}

	redirect.URL = url

	for {
		redirect.ID = helper.GenerateRandomString(model.IDLength)

		_, err = FromID(redirect.ID)
		if err == sql.ErrNoRows {
			break
		}

		if err != nil {
			return
		}
	}

	// Insert the redirect into the DB.
	_, err = db.Dot.Exec(
		db.SQL,
		"insert-redirect",
		redirect.ID,
		redirect.URL,
	)

	return
}

func InsertWithID(redirect model.Redirect) (err error) {
	_, err = FromID(redirect.ID)
	if err != nil && err != sql.ErrNoRows {
		return
	} else if err == nil {
		err = ErrIDTaken
		return
	}

	// Insert the redirect into the DB.
	_, err = db.Dot.Exec(
		db.SQL,
		"insert-redirect",
		redirect.ID,
		redirect.URL,
	)

	return
}

func scan(redirect *model.Redirect, row *sql.Row) error {
	return row.Scan(
		&redirect.ID,
		&redirect.URL,
	)
}
