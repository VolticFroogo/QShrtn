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

const (
	newResponseSuccess = iota
	newResponseInternalServerError
	newResponseForbiddenDomain
	newResponseIDTaken
	newResponseInvalidURL
)

type newReq struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

type newRes struct {
	Code  int    `json:"code"`
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// New creates a new redirect.
func New(w http.ResponseWriter, r *http.Request) {
	// Get data from the JSON request.
	var data newReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		_ = helper.JSONResponse(newRes{Code: newResponseInternalServerError, Error: err.Error()}, w)
		log.Print(err)
		return
	}

	lower := strings.ToLower(data.URL)
	host := strings.ToLower(r.Host)

	if strings.Contains(lower, host) {
		_ = helper.JSONResponse(newRes{Code: newResponseForbiddenDomain}, w)
		return
	}

	if len(data.URL) > 2048 {
		_ = helper.JSONResponse(newRes{Code: newResponseInvalidURL}, w)
		return
	}

	_, err = url.ParseRequestURI(data.URL)
	if err != nil {
		_ = helper.JSONResponse(newRes{Code: newResponseInvalidURL}, w)
		return
	}

	var redirect model.Redirect

	if data.ID == "" {
		redirect, err = Insert(data.URL)
		if err != nil {
			_ = helper.JSONResponse(newRes{Code: newResponseInternalServerError, Error: err.Error()}, w)
			log.Print(err)
			return
		}
	} else {
		redirect = model.Redirect{
			ID:  data.ID,
			URL: data.URL,
		}

		err = InsertWithShort(redirect)

		if err != nil {
			if err == ErrIDTaken {
				_ = helper.JSONResponse(newRes{Code: newResponseIDTaken}, w)
				return
			}

			_ = helper.JSONResponse(newRes{Code: newResponseInternalServerError, Error: err.Error()}, w)
			log.Print(err)
			return
		}
	}

	_ = helper.JSONResponse(newRes{
		Code: newResponseSuccess,
		ID:   redirect.ID,
	}, w)
}
