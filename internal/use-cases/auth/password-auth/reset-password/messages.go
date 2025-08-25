package resetpasswordcmd

import (
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type MessageProvider interface {
	Success() returns.Api
	Default() returns.Api
	UserNotExists() returns.Api
	UserIsAlreadyActive() returns.Api
	UserIsBlocked() returns.Api
	UserAlreadyHasCredential() returns.Api
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

func (m messages) UserIsAlreadyActive() returns.Api {
	return returns.BadRequest("user is already active.")
}

func (m messages) UserIsBlocked() returns.Api {
	return returns.Forbidden("user is block.")
}

func (m messages) UserAlreadyHasCredential() returns.Api {
	return returns.BadRequest("user already has credential.")
}
