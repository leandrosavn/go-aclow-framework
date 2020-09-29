package flow

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-aclow-framework/internal/pkg/modeltype"
	"github.com/go-aclow-framework/pkg/main/auth/model"
	"github.com/go-aclow-framework/pkg/main/config"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/google/uuid"
	"github.com/jtblin/go-ldap-client"
	"github.com/lfigueiredo82/aclow"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var (
	clientAuth = &ldap.LDAPClient{
		Base:         config.LDAPBase(),
		Host:         config.LDAPHost(),
		Port:         config.LDAPPort(),
		UseSSL:       config.LDAPUserSSL(),
		BindDN:       config.LDAPBindDN(),
		BindPassword: config.LDAPBindPassword(),
		UserFilter:   config.LDAPUserFilter(),
		GroupFilter:  config.LDAPGroupFilter(),
		Attributes:   []string{"givenName", "sn", "userPrincipalName"},
	}
)

type UserLogin struct {
	app *aclow.App

	User     string `json:"user"`
	Password string `json:"password"`
	Context  string `json:"context"`
}

func (t *UserLogin) Address() []string { return []string{"user-login"} }

func (t *UserLogin) Start(app *aclow.App) {
	t.app = app
	createUserAdmin(app)
}

func (t *UserLogin) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	login := UserLogin{}
	err := json.Unmarshal(msg.Body.([]byte), &login)
	if err != nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: err.Error(), StatusCode: 401}
	}

	if login.User == "" || login.Password == "" {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}

	client := t.app.Resources["mongo"].(*mongo.Client)
	db := client.Database(config.MongoDbDatabase())
	ctx := context.Background()

	collection := db.Collection("users")

	reply, err := call("core@load-user", aclow.Message{Body: login.User})

	var user = modeltype.User{}
	if err == nil {
		err = json.Unmarshal([]byte(reply.Body.(string)), &user)
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: err.Error(), StatusCode: 500}
		}
	}

	if user.Type == "system" {

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}

	} else {
		// It is the responsibility of the caller to close the connection
		defer clientAuth.Close()

		ok, userInfo, err := clientAuth.Authenticate(login.User, login.Password)
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}

		if !ok {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}

		filter := bson.M{
			"user": login.User,
		}

		upsert := true
		updateOptions := options.UpdateOptions{Upsert: &upsert}

		var update interface{}
		var user = modeltype.User{Profile: "user"}

		find := collection.FindOne(ctx, filter)
		if find.Err() != nil {
			update = bson.M{
				"$set": bson.M{
					"_id":       uuid.New().String(),
					"user":      login.User,
					"profile":   "user",
					"type":      "provider",
					"firstName": userInfo["givenName"],
					"lastName":  userInfo["sn"],
					"email":     userInfo["userPrincipalName"],
				},
			}
		} else {
			find.Decode(&user)
			update = bson.M{
				"$set": bson.M{
					"firstName": userInfo["givenName"],
					"lastName":  userInfo["sn"],
					"email":     userInfo["userPrincipalName"],
				},
			}
		}
		_, err = collection.UpdateOne(ctx, filter, update, &updateOptions)
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: err.Error(), StatusCode: 500}
		}
	}

	result, err := call("auth@build-token", aclow.Message{
		Body: map[string]interface{}{
			"user":           login.User,
			"context":        login.Context,
			"accountProfile": user.Profile,
			"exp":            nil,
		},
	})
	if err != nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}

	body := result.Body.(map[string]interface{})
	token := body["token"].(string)
	refresToken := body["refToken"].(string)

	tokenDetails := model.Token{
		AccessToken:          token,
		RefreshToken:         refresToken,
		Auth_Id:              login.User,
		Auth_Context:         login.Context,
		Auth_Account_Profile: user.Profile,
	}
	j, _ := json.Marshal(tokenDetails)

	return aclow.Message{Body: string(j)}, nil

}

func createUserAdmin(app *aclow.App) {

	client := app.Resources["mongo"].(*mongo.Client)
	db := client.Database(config.MongoDbDatabase())

	ctx := context.Background()

	collection := db.Collection("users")

	filter := bson.M{
		"user": "admin",
	}

	count, _ := collection.CountDocuments(ctx, filter)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		upsert := true
		updateOptions := options.UpdateOptions{Upsert: &upsert}

		update := bson.M{
			"$set": bson.M{
				"_id":      "1",
				"user":     "admin",
				"password": string(hashedPassword),
				"profile":  "admin",
				"type":     "system",
			},
		}
		_, err = collection.UpdateOne(ctx, filter, update, &updateOptions)
		if err != nil {
			log.Fatal(err)
		}
	}

}
