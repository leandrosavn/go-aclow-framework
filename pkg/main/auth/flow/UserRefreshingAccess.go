package flow

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-aclow-framework/pkg/main/auth/config"
	"github.com/go-aclow-framework/pkg/main/auth/model"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/lfigueiredo82/aclow"
)

type UserRefreshToken struct {
}

func (t *UserRefreshToken) Address() []string { return []string{"user-refreshing-access"} }

func (t *UserRefreshToken) Start(app *aclow.App) {}

func (t *UserRefreshToken) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	tokenDetalis := model.Token{}
	err := json.Unmarshal(msg.Body.([]byte), &tokenDetalis)
	if err != nil {
		log.Println(err)
	}

	//verify the token
	token, err := jwt.Parse(tokenDetalis.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.RefreshSecret()), nil
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
	if ok && token.Valid {
		//check last token
		if tokenDetalis.AccessToken != claims["token"].(string) {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}
		user := claims["auth_id"].(string)
		context := claims["auth_context"].(string)
		accountProfile := claims["auth_account_profile"].(string)
		exp := claims["exp"]

		result, err := call("auth@build-token", aclow.Message{
			Body: map[string]interface{}{
				"user":           user,
				"context":        context,
				"accountProfile": accountProfile,
				"exp":            exp,
			},
		})
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}

		body := result.Body.(map[string]interface{})
		token := body["token"].(string)
		refToken := body["refToken"].(string)

		tokenDetails := model.Token{
			AccessToken:          token,
			RefreshToken:         refToken,
			Auth_Id:              user,
			Auth_Context:         context,
			Auth_Account_Profile: accountProfile,
		}

		j, _ := json.Marshal(tokenDetails)

		return aclow.Message{Body: string(j)}, nil
	} else {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}

}
