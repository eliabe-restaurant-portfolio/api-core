package createusercmd

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type MessageProvider interface {
	Success(created *entities.User) returns.Api
	Default() returns.Api
	RepeatedUser() returns.Api
}

type messages struct{}

func NewMessages() MessageProvider {
	return messages{}
}

func (m messages) Success(created *entities.User) returns.Api {
	return returns.Success("user created", &Return{
		UserToken: created.Token.String(),
		Status:    string(created.Status),
	})
}

func (m messages) Default() returns.Api {
	return returns.InternalServerError([]string{})
}

func (m messages) RepeatedUser() returns.Api {
	return returns.BadRequest("repeated user")
}
