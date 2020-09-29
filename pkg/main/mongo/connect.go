package mongo

import (
	"github.com/go-aclow-framework/pkg/main/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDbDsn()))

	return client, err

}
