package flow

import (
	"encoding/json"

	"github.com/go-aclow-framework/pkg/main/api/config"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/lfigueiredo82/aclow"
)

type Authenticating struct {
	Context     string `json:"context"`
	application *aclow.App
	routeConfig *config.Config
}

func (t *Authenticating) Address() []string { return []string{"authenticating"} }

func (t *Authenticating) Start(app *aclow.App) {
	t.application = app
	t.routeConfig = config.ConfigRoutes()
}

func (t *Authenticating) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	auth := Authenticating{}
	err := json.Unmarshal(msg.Body.([]byte), &auth)
	if err != nil {
		return aclow.Message{Body: ""}, &utils.MessageError{ErrorMessage: "No AUTH parameter", StatusCode: 500}
	}

	context := auth.Context
	authFlow := t.routeConfig.Context[context]["login-event"]
	if authFlow != "" {
		reply, err := t.application.Call(authFlow, msg)
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}
		if reply.Body != nil {
			return aclow.Message{Body: reply.Body}, nil
		} else {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
		}

	} else {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}
}
