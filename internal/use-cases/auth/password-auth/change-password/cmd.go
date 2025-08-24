package changepasswordcmd

import (
	"context"

	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/eliabe-portfolio/restaurant-app/internal/aggregates"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	userrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/users"
	uow "github.com/eliabe-portfolio/restaurant-app/internal/unit-of-work"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	hashing "github.com/eliabe-portfolio/restaurant-app/pkg/hash"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/google/uuid"
)

type Command struct {
	messages       MessageProvider
	unitOfWork     uow.UnitOfWork
	userRepository userrepo.UserRepository
}

type Params struct {
	ActorToken  uuid.UUID
	Email       valueobjects.Email
	OldPassword valueobjects.Password
	NewPassword valueobjects.Password
}

type RelatedEntities struct {
	Actor *entities.User
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

	actorEntity, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Token: params.ActorToken,
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	var actor = aggregates.NewUser(actorEntity)

	if !actor.Exists() {
		return cmd.messages.UserNotExists(), nil
	}

	if actor.IsInactive() {
		return cmd.messages.UserIsInactive(), nil
	}

	if actor.IsBlocked() {
		return cmd.messages.UserIsBlocked(), nil
	}

	isValid := actor.IsValidPassword(params.OldPassword)

	if !isValid && actor.HasMaxFailedLoginAttempts() {
		err := cmd.userRepository.Update(userrepo.UpdateUserDto{
			Ctx:       ctx,
			UserToken: actorEntity.Token,
			Status:    &constants.UserBloqued,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}

		cmd.unitOfWork.Commit(ctx)

		return cmd.messages.BlockUser(), nil
	}

	if !isValid && !actor.HasMaxFailedLoginAttempts() {
		actorEntity.FailedLoginAttempts++

		err = cmd.userRepository.Update(userrepo.UpdateUserDto{
			Ctx:                 ctx,
			UserToken:           actorEntity.Token,
			FailedLoginAttempts: &actorEntity.FailedLoginAttempts,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}

		cmd.unitOfWork.Commit(ctx)

		return cmd.messages.UpdateUserFailedLoginAttemps(actor), err
	}

	if isValid && actor.HasMaxFailedLoginAttempts() {
		actorEntity.FailedLoginAttempts = 0

		err = cmd.userRepository.Update(userrepo.UpdateUserDto{
			Ctx:                 ctx,
			UserToken:           actorEntity.Token,
			FailedLoginAttempts: &actorEntity.FailedLoginAttempts,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}
	}

	hash, err := hashing.Hash(params.NewPassword.Get())
	if err != nil {
		return cmd.messages.Default(), err
	}

	if err := cmd.userRepository.Update(userrepo.UpdateUserDto{
		Ctx:       ctx,
		UserToken: params.ActorToken,
		Password:  hash,
	}); err != nil {
		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(ctx)

	return cmd.messages.Success(), nil
}
