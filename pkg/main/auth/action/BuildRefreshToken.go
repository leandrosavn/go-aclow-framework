package action

import (
	"time"

	"github.com/go-aclow-framework/pkg/main/auth/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/lfigueiredo82/aclow"
)

type BuildRefreshToken struct {
}

func (t *BuildRefreshToken) Address() []string { return []string{"build-refresh-token"} }

func (t *BuildRefreshToken) Start(app *aclow.App) {}

func (t *BuildRefreshToken) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	params := msg.Body.(map[string]interface{})
	user := params["user"].(string)
	context := params["context"].(string)
	accountProfile := params["accountProfile"].(string)
	token := params["token"].(string)
	var exp interface{} = nil
	if params["exp"] != nil {
		exp = params["exp"].(interface{})
	}

	//
	var err error
	//Creating Access Token

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["auth_id"] = user
	atClaims["auth_context"] = context
	atClaims["auth_account_profile"] = accountProfile
	atClaims["token"] = token
	if exp == nil {
		atClaims["exp"] = time.Now().Add(time.Hour * 24 * 5).Unix()
	} else {
		atClaims["exp"] = exp
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	refToken, err := at.SignedString([]byte(config.RefreshSecret()))
	if err != nil {
		return aclow.Message{Body: ""}, err
	}
	return aclow.Message{Body: refToken}, nil

}
