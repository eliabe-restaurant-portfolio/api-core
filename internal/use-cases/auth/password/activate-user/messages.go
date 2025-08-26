package activateusercmd

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type MessageProvider interface {
	Success() returns.Api
	Default() returns.Api
	ResetPasswordNotExists() returns.Api
	UserIsAlreadyActive() returns.Api
	UserIsBlocked() returns.Api
	InvalidResetToken() returns.Api
}

type messages struct{}

func NewMessages() MessageProvider {
	return messages{}
}

func (m messages) Success() returns.Api {
	return returns.Success("user activated", Return{UserStatus: string(constants.UserActive)})
}

func (m messages) Default() returns.Api {
	return returns.InternalServerError([]string{})
}

func (m messages) ResetPasswordNotExists() returns.Api {
	return returns.BadRequest("reset password not exists")
}

func (m messages) UserIsAlreadyActive() returns.Api {
	return returns.BadRequest("user is already active.")
}

func (m messages) UserIsBlocked() returns.Api {
	return returns.Forbidden("user is blocked.")
}

func (m messages) InvalidResetToken() returns.Api {
	return returns.BadRequest("token invalid.")
}
