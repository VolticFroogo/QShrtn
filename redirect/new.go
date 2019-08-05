package redirect

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/VolticFroogo/QShrtn/helper"
	"github.com/VolticFroogo/QShrtn/model"
)

type newReq struct {
	URL string
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

	url := strings.ToLower(data.URL)

	if strings.Contains(url, "qshr.tn") {
		helper.JSONResponse(model.Code{Code: model.ResponseForbiddenDomain}, w)
		return
	}

	redirect, err := Insert(url)
	if err != nil {
		helper.JSONResponse(model.Code{Code: model.ResponseInternalServerError}, w)
		log.Print(err)
		return
	}

	helper.JSONResponse(newRes{
		Code: model.ResponseSuccess,
		ID:   redirect.ID,
	}, w)
}
