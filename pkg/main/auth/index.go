package auth

import (
	"github.com/go-aclow-framework/pkg/main/auth/action"
	"github.com/go-aclow-framework/pkg/main/auth/flow"
	"github.com/lfigueiredo82/aclow"
)

func Nodes() []aclow.Node {
	return []aclow.Node{
		&flow.UserLogin{},
		&flow.UserRefreshToken{},
		&flow.UserValidatingAccess{},

		&action.BuildToken{},
		&action.BuildRefreshToken{},
	}
}
