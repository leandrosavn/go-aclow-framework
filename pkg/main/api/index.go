package api

import (
	"github.com/go-aclow-framework/pkg/main/api/flow"
	"github.com/lfigueiredo82/aclow"
)

func Nodes() []aclow.Node {
	return []aclow.Node{
		&flow.Ping{},
		&flow.CallingFlow{},
		&flow.Authenticating{},
	}
}
