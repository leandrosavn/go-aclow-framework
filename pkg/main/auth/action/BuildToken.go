package action

import (
	"time"

	"github.com/go-aclow-framework/pkg/main/auth/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/lfigueiredo82/aclow"
)

type BuildToken struct {
}

func (t *BuildToken) Address() []string { return []string{"build-token"} }

func (t *BuildToken) Start(app *aclow.App) {}

func (t *BuildToken) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	params := msg.Body.(map[string]interface{})
	user := params["user"].(string)
	context := params["context"].(string)
	accountProfile := params["accountProfile"].(string)

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["auth_id"] = user
	atClaims["auth_context"] = context
	atClaims["auth_account_type"] = accountProfile
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.AccessSecret()))
	if err != nil {
		return aclow.Message{Body: map[string]interface{}{
			"token":    "",
			"refToken": ""},
		}, err
	}

	params["token"] = token
	refToken, err := call("auth@build-refresh-token", aclow.Message{Body: params})
	if err != nil {
		return aclow.Message{
			Body: map[string]interface{}{
				"token":    "",
				"refToken": ""},
		}, err
	}
	return aclow.Message{
		Body: map[string]interface{}{
			"token":    token,
			"refToken": refToken.Body.(string)},
	}, nil
}
