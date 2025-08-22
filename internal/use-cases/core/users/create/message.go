package createusercmd

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
)

type MessageProvider interface {
	Success(created *entities.User) returns.Api
	Default() returns.Api
	ExistsUserWithSameEmail() returns.Api
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

func (m messages) ExistsUserWithSameEmail() returns.Api {
	return returns.BadRequest("exists user with same email")
}
