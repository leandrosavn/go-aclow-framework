package mongo

import (
	"github.com/go-aclow-framework/pkg/main/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(env.MONGO_URI))

	return client, err

}
