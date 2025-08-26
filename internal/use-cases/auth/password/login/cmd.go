package traditionallogincmd

import (
	"context"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/aggregates"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/envs"
	userrepo "github.com/eliabe-restaurant-portfolio/api-core/internal/repositories/users"
	uow "github.com/eliabe-restaurant-portfolio/api-core/internal/unit-of-work"
	objects "github.com/eliabe-restaurant-portfolio/api-core/internal/value-objects"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/jwt"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type Command struct {
	messages       MessageProvider
	unitOfWork     uow.UnitOfWork
	userRepository userrepo.UserRepository
}

type RelatedEntities struct {
	Actor *entities.User
}

type Params struct {
	Email    objects.Email
	Password objects.Password
}

type Return struct {
	UserToken   string `json:"user_token"`
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
	IssuedAt    string `json:"issued_at"`
}

func New(adapters *adapters.Adapters) Command {
	return Command{
		messages:       NewMessages(),
		unitOfWork:     (*adapters).UnitOfWork(),
		userRepository: (*adapters).Repositories().User(),
	}
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	var err error
	var ctx = context.Background()
	cmd.unitOfWork.Init(ctx)

	related, err := cmd.getRelatedEntities(params)
	if err != nil {
		return cmd.messages.Default(), err
	}

	var actor = aggregates.NewUser(related.Actor)

	if !actor.Exists() {
		return cmd.messages.UserNotExists(), nil
	}

	if actor.IsInactive() {
		return cmd.messages.UserIsInactive(), nil
	}

	if actor.IsBlocked() {
		return cmd.messages.UserIsBlocked(), nil
	}

	isValid := actor.IsValidPassword(params.Password)

	if !isValid && actor.HasMaxFailedLoginAttempts() {
		err := cmd.userRepository.Update(userrepo.UpdateUserDto{
			Ctx:       ctx,
			UserToken: related.Actor.Token,
			Status:    &constants.UserBloqued,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}

		cmd.unitOfWork.Commit(ctx)

		return cmd.messages.BlockUser(), nil
	}

	if !isValid && !actor.HasMaxFailedLoginAttempts() {
		related.Actor.FailedLoginAttempts++

		err = cmd.userRepository.Update(userrepo.UpdateUserDto{
			Ctx:                 ctx,
			UserToken:           related.Actor.Token,
			FailedLoginAttempts: &related.Actor.FailedLoginAttempts,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}

		cmd.unitOfWork.Commit(ctx)

		return cmd.messages.UpdateUserFailedLoginAttemps(actor), err
	}

	if isValid && actor.HasMaxFailedLoginAttempts() {
		related.Actor.FailedLoginAttempts = 0

		err = cmd.userRepository.Update(userrepo.UpdateUserDto{
			Ctx:                 ctx,
			UserToken:           related.Actor.Token,
			FailedLoginAttempts: &related.Actor.FailedLoginAttempts,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}
	}

	accessDetails, err := jwt.Create(jwt.JwtCreateInput{
		Content:        related.Actor.Token.String(),
		Duration:       constants.AuthAccessDurantion,
		AccessClientId: envs.Get(envs.ACCESS_CLIENT_ID),
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(ctx)

	return cmd.messages.Success(related.Actor, accessDetails), nil
}

func (cmd Command) getRelatedEntities(params Params) (*RelatedEntities, error) {
	user, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Email: &params.Email,
	})
	if err != nil {
		return nil, err
	}

	return &RelatedEntities{
		Actor: user,
	}, nil
}
