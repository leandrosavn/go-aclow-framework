package core

import (
	"github.com/go-aclow-framework/pkg/main/core/action"
	"github.com/go-aclow-framework/pkg/main/core/flow"
	"github.com/lfigueiredo82/aclow"
)

func Nodes() []aclow.Node {
	return []aclow.Node{
		// flows

		&flow.LoadingUser{},

		// actions
		&action.ChangeUserType{},
		&action.LoadUser{},
	}
}
