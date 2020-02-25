package redirect

import (
	"context"
	"fmt"
	"strings"

	"github.com/VolticFroogo/QShrtn/db"
	"github.com/VolticFroogo/QShrtn/helper"
	"github.com/VolticFroogo/QShrtn/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrIDTaken = fmt.Errorf("id already taken")
)

// FromShort gets a redirect from the database given a ID.
func FromShort(id string) (redirect model.Redirect, err error) {
	ctx := context.Background()

	err = db.Redirect.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&redirect)

	return
}

// FromLong gets a redirect from the database given a URL.
func FromURL(url string) (redirect model.Redirect, err error) {
	ctx := context.Background()

	err = db.Redirect.FindOne(ctx, bson.M{
		"url": url,
	}).Decode(&redirect)

	return
}

// Insert a redirect into the database.
func Insert(url string) (redirect model.Redirect, err error) {
	ctx := context.Background()

	// Check if there's already a redirect for that URL.
	redirect, err = FromURL(url)
	// If there is an entry for this URL already, return that.
	// Or, if there was an error that wasn't no documents, return that.
	if err == nil || err != mongo.ErrNoDocuments {
		return
	}

	redirect.URL = url

	for {
		redirect.ID = helper.GenerateRandomString(model.IDLength)

		count, err := db.Redirect.CountDocuments(ctx, bson.M{
			"_id": redirect.ID,
		})
		if err != nil {
			return redirect, err
		}

		if count == 0 {
			break
		}
	}

	// Insert the redirect into the DB.
	_, err = db.Redirect.InsertOne(ctx, redirect)
	return
}

func InsertWithShort(redirect model.Redirect) (err error) {
	ctx := context.Background()

	// Check if this ID is already in use.
	// Note: this isn't necessary as unique collisions will be detected below,
	// however this will improve performance on secondary DB machines.
	// Secondaries can count at very low cost using indexes, but writes require network calls.
	count, err := db.Redirect.CountDocuments(ctx, bson.M{
		"_id": redirect.ID,
	})
	if err != nil {
		return
	}

	if count != 0 {
		err = ErrIDTaken
		return
	}

	_, err = db.Redirect.InsertOne(ctx, redirect)

	// Check if the error is a unique key collision.
	if err != nil && strings.Contains(err.Error(), "duplicate key error") {
		err = ErrIDTaken
	}

	return
}
