package traditionallogincmd

import (
	"fmt"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/aggregates"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/jwt"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type MessageProvider interface {
	Success(actor *entities.User, accessDetails *jwt.AccessDetails) returns.Api
	Default() returns.Api
	UserNotExists() returns.Api
	UserIsInactive() returns.Api
	UserIsBlocked() returns.Api
	BlockUser() returns.Api
	UpdateUserFailedLoginAttemps(actor aggregates.User) returns.Api
}

type messages struct{}

func NewMessages() MessageProvider {
	return messages{}
}

func (m messages) Success(actor *entities.User, accessDetails *jwt.AccessDetails) returns.Api {
	res := Return{
		UserToken:   actor.Token.String(),
		Status:      string(actor.Status),
		AccessToken: accessDetails.Token,
		ExpiresAt:   accessDetails.ExpiresAt,
		IssuedAt:    accessDetails.IssuedAt,
	}
	return returns.Success("login successfully", &res)
}

func (m messages) Default() returns.Api {
	return returns.InternalServerError([]string{})
}

func (m messages) UserNotExists() returns.Api {
	return returns.NotFound("user not found.")
}

func (m messages) UserIsInactive() returns.Api {
	return returns.Unauthorized("user is inactive.")
}

func (m messages) UserIsBlocked() returns.Api {
	return returns.Unauthorized("user is block.")
}

func (m messages) BlockUser() returns.Api {
	return returns.Success("user has blocked.", nil)
}

func (m messages) UpdateUserFailedLoginAttemps(actor aggregates.User) returns.Api {
	message := fmt.Sprintf("invalid credentials. failed attempt: %v", actor.GetFailedLoginAttempts())
	return returns.Unauthorized(message)
}
