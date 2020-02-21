package redirect

import (
	"context"
	"fmt"
	"log"
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

		_, err = FromShort(redirect.ID)
		if err == mongo.ErrNoDocuments {
			break
		}

		if err != nil {
			return
		}
	}

	// Insert the redirect into the DB.
	_, err = db.Redirect.InsertOne(ctx, redirect)
	return
}

func InsertWithShort(redirect model.Redirect) (err error) {
	ctx := context.Background()

	_, err = db.Redirect.InsertOne(ctx, redirect)

	// Check if the error is a unique key collision.
	if err != nil && strings.Contains(err.Error(), "duplicate key error") {
		log.Print(err)
		err = ErrIDTaken
	}

	return
}
