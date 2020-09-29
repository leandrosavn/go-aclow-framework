package flow

import (
	"github.com/lfigueiredo82/aclow"
)

type Ping struct {
}

func (t *Ping) Address() []string { return []string{"ping"} }

func (t *Ping) Start(app *aclow.App) {}

func (t *Ping) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {
	return aclow.Message{Body: "pong"}, nil
}
