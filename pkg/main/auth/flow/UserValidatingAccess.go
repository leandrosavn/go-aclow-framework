package flow

import (
	"fmt"

	"github.com/go-aclow-framework/pkg/main/auth/config"
	"github.com/go-aclow-framework/pkg/main/auth/model"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/lfigueiredo82/aclow"
)

type UserValidatingAccess struct {
}

func (t *UserValidatingAccess) Address() []string { return []string{"user-validating-access"} }

func (t *UserValidatingAccess) Start(app *aclow.App) {}

func (t *UserValidatingAccess) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	header := msg.Header
	userToken := header["token"].(string)

	//verify the token
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AccessSecret()), nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok || !token.Valid {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}

	user := claims["auth_id"].(string)
	context := claims["auth_context"].(string)
	userValidatingAccess := model.UserAccess{AuthId: user, AuthContext: context}

	return aclow.Message{Body: userValidatingAccess}, nil

}
