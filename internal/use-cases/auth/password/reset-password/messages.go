package resetpasswordcmd

import (
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type MessageProvider interface {
	Success() returns.Api
	Default() returns.Api
	UserNotExists() returns.Api
}

type messages struct{}

func NewMessages() MessageProvider {
	return messages{}
}

func (m messages) Success() returns.Api {
	return returns.Success("reset password successufully", nil)
}

func (m messages) Default() returns.Api {
	return returns.InternalServerError([]string{})
}

func (m messages) UserNotExists() returns.Api {
	return returns.NotFound("user not found.")
}
