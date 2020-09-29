package action

import (
	"context"
	"encoding/json"

	"github.com/go-aclow-framework/internal/pkg/modeltype"
	"github.com/go-aclow-framework/pkg/main/config"
	"github.com/go-aclow-framework/pkg/main/core/model"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/lfigueiredo82/aclow"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChangeUserType struct {
	app *aclow.App
}

func (t *ChangeUserType) Address() []string { return []string{"change-user-type"} }

func (t *ChangeUserType) Start(app *aclow.App) {
	t.app = app
}

func (t *ChangeUserType) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	user := msg.Body.(model.UserType)

	client := t.app.Resources["mongo"].(*mongo.Client)
	db := client.Database(config.MongoDbDatabase())

	ctx := context.Background()

	collection := db.Collection("users")

	var filter = bson.M{
		"user": user.User,
	}

	update := bson.M{
		"$set": bson.M{
			"profile": user.Profile,
		},
	}

	returnDocument := options.After
	updateOptions := options.FindOneAndUpdateOptions{ReturnDocument: &returnDocument}
	res := collection.FindOneAndUpdate(ctx, filter, update, &updateOptions)
	if res == nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Documento nao encontrado", StatusCode: 500}
	}

	usuario := modeltype.User{}
	res.Decode(&usuario)
	j, _ := json.Marshal(usuario)
	return aclow.Message{Body: string(j)}, nil
}
