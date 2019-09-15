package redirect

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/VolticFroogo/QShrtn/helper"
	"github.com/VolticFroogo/QShrtn/model"
)

type newReq struct {
	URL, ID string
}

type newRes struct {
	Code int
	ID   string
}

// New creates a new redirect.
func New(w http.ResponseWriter, r *http.Request) {
	// Get data from the JSON request.
	var data newReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helper.JSONResponse(model.Code{Code: model.ResponseInternalServerError}, w)
		log.Print(err)
		return
	}

	lower := strings.ToLower(data.URL)

	if strings.Contains(lower, "qshr.tn") {
		helper.JSONResponse(model.Code{Code: model.ResponseForbiddenDomain}, w)
		return
	}

	_, err = url.ParseRequestURI(data.URL)
	if err != nil {
		helper.JSONResponse(model.Code{Code: model.ResponseInvalidURL}, w)
		return
	}

	var redirect model.Redirect

	if data.ID == "" {
		redirect, err = Insert(data.URL)
		if err != nil {
			helper.JSONResponse(model.Code{Code: model.ResponseInternalServerError}, w)
			log.Print(err)
			return
		}
	} else {
		redirect = model.Redirect{
			ID:  data.ID,
			URL: data.URL,
		}

		err = InsertWithID(redirect)

		if err != nil {
			if err == ErrIDTaken {
				helper.JSONResponse(model.Code{Code: model.ResponseIDTaken}, w)
				return
			}

			helper.JSONResponse(model.Code{Code: model.ResponseInternalServerError}, w)
			log.Print(err)
			return
		}
	}

	helper.JSONResponse(newRes{
		Code: model.ResponseSuccess,
		ID:   redirect.ID,
	}, w)
}
