package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	uri = os.Getenv("DB")

	Redirect *mongo.Collection
)

// Init initialises the database.
func Init() (err error) {
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(opts)
	if err != nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return
	}

	err = client.Ping(ctx, readpref.Nearest())
	if err != nil {
		return
	}

	db := client.Database("qshrtn")

	Redirect = db.Collection("redirect")
	return
}
