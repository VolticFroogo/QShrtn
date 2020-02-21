package db

import (
	"context"
	"time"

	"github.com/VolticFroogo/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	configDirectory = "configs/db.ini"
)

var (
	Redirect *mongo.Collection
)

// Config is the config structure.
type Config struct {
	URI string
}

// Init initialises the database.
func Init() (err error) {
	// Load the config.
	cfg := Config{}
	err = config.Load(configDirectory, &cfg)
	if err != nil {
		return
	}

	opts := options.Client().ApplyURI(cfg.URI)

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
