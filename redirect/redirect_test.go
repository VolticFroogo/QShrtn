package redirect

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/VolticFroogo/QShrtn/db"
	"github.com/VolticFroogo/QShrtn/helper"
	"github.com/VolticFroogo/QShrtn/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	contentType = "application/json; charset=UTF-8"
)

const (
	urlPrefix = "http://localhost:8080/"
)

var (
	// Helper variables for testing.
	ctx    = context.Background()
	client = http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func init() {
	err := db.Init()
	if err != nil {
		log.Fatalf("Could not initialise database: %v", err)
	}

	helper.Seed()

	go func() {
		// Create a new Mux Router with strict slash.
		r := mux.NewRouter()
		r.StrictSlash(true)

		// Handle new URL requests.
		r.Handle("/new/", http.HandlerFunc(New)).Methods(http.MethodPost)

		// Handle all unknown links, possibly redirecting links.
		r.Handle("/{id}", http.HandlerFunc(Handle))
		r.Handle("/{id}/json", http.HandlerFunc(JSON))

		log.Print("Listening for incoming HTTP requests on port 8080.")

		// Serve plain HTTP responses.
		err = http.ListenAndServe(":8080", r)

		if err != nil {
			log.Fatalf("Could not start handler: %v", err)
		}
	}()
}

func TestHandle(t *testing.T) {
	assert := assert.New(t)

	// Insert a test entry into the database.
	_, err := db.Redirect.InsertOne(ctx, model.Redirect{
		ID:  "test",
		URL: "https://froogo.co.uk/",
	})
	assert.Nil(err)

	t.Log("Checking an existing link.")
	res, err := client.Get(urlPrefix + "test")
	assert.Nil(err)
	assert.Equal(http.StatusMovedPermanently, res.StatusCode)
	assert.Equal("https://froogo.co.uk/", res.Header.Get("location"))

	t.Log("Checking a non-existent link.")
	res, err = client.Get(urlPrefix + "unknown")
	assert.Nil(err)
	assert.Equal(http.StatusTemporaryRedirect, res.StatusCode)
	assert.Equal("/not-found/", res.Header.Get("location"))
}

func TestJSON(t *testing.T) {
	assert := assert.New(t)

	// Insert a test entry into the database.
	_, err := db.Redirect.InsertOne(ctx, model.Redirect{
		ID:  "json",
		URL: "https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON",
	})
	assert.Nil(err)

	var decoded jsonResponse

	t.Log("Checking an existing link via the JSON API.")
	res, err := client.Get(urlPrefix + "json/json")
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&decoded)
	assert.Nil(err)
	assert.Equal(responseSuccess, decoded.Code)
	assert.Equal("https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON", decoded.URL)

	t.Log("Checking a non-existent link via the JSON API.")
	res, err = client.Get(urlPrefix + "unknown/json")
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&decoded)
	assert.Nil(err)
	assert.Equal(responseNotFound, decoded.Code)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	t.Log("Inserting a valid redirect with unspecified ID.")
	_, decoded, document := insert(newReq{
		URL: "https://atrello.co.uk/",
	}, assert, true)
	assert.Equal(newResponseSuccess, decoded.Code)
	assert.Equal("https://atrello.co.uk/", document.URL)

	t.Log("Inserting a valid redirect with specified ID.")
	_, decoded, document = insert(newReq{
		URL: "https://duckduckgo.com/",
		ID:  "ddg",
	}, assert, true)
	assert.Equal(newResponseSuccess, decoded.Code)
	assert.Equal("https://duckduckgo.com/", document.URL)
	assert.Equal("ddg", document.ID)

	t.Log("Inserting an invalid redirect due to the ID already being taken.")
	_, err := db.Redirect.InsertOne(ctx, model.Redirect{
		ID:  "taken",
		URL: "https://github.com/VolticFroogo/",
	})
	assert.Nil(err)
	_, decoded, document = insert(newReq{
		URL: "https://stackoverflow.com/",
		ID:  "taken",
	}, assert, false)
	assert.Equal(newResponseIDTaken, decoded.Code)

	t.Log("Inserting an invalid redirect due to the URL being invalid.")
	_, decoded, _ = insert(newReq{
		URL: "invalid-url",
	}, assert, false)
	assert.Equal(newResponseInvalidURL, decoded.Code)

	t.Log("Inserting an invalid redirect due to the URL being too long.")
	_, decoded, _ = insert(newReq{
		URL: "https://froogo.co.uk/" + helper.GenerateRandomString(model.MaxURLLength),
	}, assert, false)
	assert.Equal(newResponseInvalidURL, decoded.Code)

	t.Log("Inserting an invalid redirect due to the URL containing a forbidden domain.")
	_, decoded, _ = insert(newReq{
		URL: "https://" + domain + "/test",
	}, assert, false)
	assert.Equal(newResponseForbiddenDomain, decoded.Code)
}

func insert(request newReq, assert *assert.Assertions, valid bool) (res *http.Response, decoded newRes, document model.Redirect) {
	req, err := json.Marshal(request)
	assert.Nil(err)

	res, err = client.Post(urlPrefix+"new/", contentType, bytes.NewBuffer(req))
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&decoded)
	assert.Nil(err)

	if valid {
		err = db.Redirect.FindOne(ctx, bson.M{
			"_id": decoded.ID,
		}).Decode(&document)
		assert.Nil(err)
	}

	return
}
