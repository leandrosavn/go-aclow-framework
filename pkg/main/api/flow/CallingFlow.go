package flow

import (
	"github.com/go-aclow-framework/pkg/main/api/config"
	"github.com/go-aclow-framework/pkg/main/auth/model"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/lfigueiredo82/aclow"
)

type CallingFlow struct {
	application *aclow.App
	routeConfig *config.Config
}

func (t *CallingFlow) Address() []string { return []string{"calling-flow"} }

func (t *CallingFlow) Start(app *aclow.App) {
	t.application = app
	t.routeConfig = config.ConfigRoutes()
}

func (t *CallingFlow) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {

	header := msg.Header
	address := header["address"].(string)

	event := t.routeConfig.Event[address]
	if event != nil {
		contextId := event["auth"]
		context := t.routeConfig.Context[contextId]
		if context != nil {
			return check(t.application, msg, address, contextId, context)
		} else {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Not found", StatusCode: 404}
		}
	} else {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Not found", StatusCode: 404}
	}

}

func check(application *aclow.App, msg aclow.Message, address string, contextId string, context map[string]string) (aclow.Message, error) {

	checkEvent := context["check-token-event"]
	if checkEvent != "" {
		tokenResult, err := application.Call(checkEvent, aclow.Message{Header: msg.Header, Body: msg.Body})
		if err != nil {
			return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Check Token Error", StatusCode: 500}
		}
		userValidatingAccess := tokenResult.Body.(model.UserAccess)
		return callEvent(application, userValidatingAccess, contextId, context, address, msg)
	} else {
		userValidatingAccess := model.UserAccess{AuthContext: contextId, AuthId: contextId}
		return callEvent(application, userValidatingAccess, contextId, context, address, msg)
	}

}

func callEvent(application *aclow.App, tokenResult model.UserAccess, contextId string, context map[string]string, address string, msg aclow.Message) (aclow.Message, error) {

	headers := msg.Header

	headers["auth_id"] = tokenResult.AuthId
	headers["auth_context"] = tokenResult.AuthContext

	if tokenResult.AuthId != "" && contextId == tokenResult.AuthContext {
		reply, err := application.Call(address, aclow.Message{Header: headers, Body: msg.Body})
		if err != nil {
			if err, ok := err.(*utils.MessageError); ok {
				return aclow.Message{Body: nil}, err
			}

		}
		return aclow.Message{Body: reply.Body}, nil
	} else {
		return aclow.Message{Body: nil}, &utils.MessageError{ErrorMessage: "Unauthorized", StatusCode: 401}
	}

}
