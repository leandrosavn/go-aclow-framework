package action

import (
	"context"
	"encoding/json"

	"github.com/go-aclow-framework/internal/pkg/modeltype"
	"github.com/go-aclow-framework/pkg/main/config"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/lfigueiredo82/aclow"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoadUser struct {
	app *aclow.App
}

func (t *LoadUser) Address() []string { return []string{"load-user"} }

func (t *LoadUser) Start(app *aclow.App) {
	t.app = app
}

func (t *LoadUser) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	paramUser := msg.Body.(string)

	client := t.app.Resources["mongo"].(*mongo.Client)
	db := client.Database(config.MongoDbDatabase())

	ctx := context.Background()

	collection := db.Collection("users")

	user := modeltype.User{}
	err := collection.FindOne(ctx, bson.M{"user": paramUser}).Decode(&user)
	if err != nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: err.Error(), StatusCode: 500}
	}
	j, _ := json.Marshal(user)
	return aclow.Message{Body: string(j)}, nil

}
