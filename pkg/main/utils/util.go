package utils

import (
	"context"
	"log"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CursorToArray(ctx context.Context, cursor *mongo.Cursor) bson.A {
	var result = bson.A{}

	for cursor.Next(ctx) {
		elem := &bson.M{}
		if err := cursor.Decode(elem); err != nil {
			log.Println(err)
		} else {
			result = append(result, *elem)
		}
	}

	return result
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func NameOf(key string) string {
	r, _ := regexp.Compile(".*\\/")
	return r.ReplaceAllString(key, "")
}
