package deleteusercmd

import (
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
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
	return returns.Success("user deleted", nil)
}

func (m messages) Default() returns.Api {
	return returns.InternalServerError([]string{})
}

func (m messages) UserNotExists() returns.Api {
	return returns.BadRequest("user not exists")
}
