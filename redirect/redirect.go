package redirect

import (
	"database/sql"
	"net/http"

	"github.com/VolticFroogo/QShrtn/helper"
	"github.com/gorilla/mux"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	// Get the variables from the URL.
	vars := mux.Vars(r)

	id := vars["id"]

	redirect, err := FromID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/not-found/", http.StatusTemporaryRedirect)
			return
		}

		helper.ThrowErr(w, r, err)
		return
	}

	// Let clients cache our response for speed.
	// Public: responses will always be the same for all clients, proxies can cache this.
	// Max-Age: cache this response for one year.
	// Immutable: the response will never change.
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}
