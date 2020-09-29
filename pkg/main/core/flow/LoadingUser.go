package flow

import (
	"encoding/json"

	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/lfigueiredo82/aclow"
)

type LoadingUser struct {
	User string `json:"user"`
}

func (t *LoadingUser) Address() []string { return []string{"loading-user"} }

func (t *LoadingUser) Start(app *aclow.App) {}

func (t *LoadingUser) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	userParam := LoadingUser{}
	err := json.Unmarshal(msg.Body.([]byte), &userParam)
	if err != nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: err.Error(), StatusCode: 500}
	}
	result, err := call("core@load-user", aclow.Message{Body: userParam.User})
	if err != nil {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: err.Error(), StatusCode: 500}
	}

	return aclow.Message{Body: result.Body.(string)}, nil

}
