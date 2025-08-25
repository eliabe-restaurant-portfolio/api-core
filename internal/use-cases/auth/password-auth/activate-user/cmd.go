package activateusercmd

import (
	"context"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/aggregates"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	resetpasswordrepo "github.com/eliabe-restaurant-portfolio/api-core/internal/repositories/reset-password"
	userrepo "github.com/eliabe-restaurant-portfolio/api-core/internal/repositories/users"
	uow "github.com/eliabe-restaurant-portfolio/api-core/internal/unit-of-work"
	valueobjects "github.com/eliabe-restaurant-portfolio/api-core/internal/value-objects"
	hashing "github.com/eliabe-restaurant-portfolio/api-core/pkg/hash"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
	"github.com/google/uuid"
)

type Command struct {
	messages                MessageProvider
	unitOfWork              uow.UnitOfWork
	userRepository          userrepo.UserRepository
	resetPasswordRepository resetpasswordrepo.ResetPasswordRepository
}

type Params struct {
	ResetPasswordToken uuid.UUID
	ResetPasswordHash  string
	NewPassword        valueobjects.Password
}

type Return struct {
	UserStatus string `json:"user_status"`
}

func New(adapters *adapters.Adapters) Command {
	return Command{
		messages:                NewMessages(),
		unitOfWork:              (*adapters).UnitOfWork(),
		userRepository:          (*adapters).Repositories().User(),
		resetPasswordRepository: (*adapters).Repositories().ResetPassword(),
	}
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	var err error
	var ctx = context.Background()
	cmd.unitOfWork.Init(ctx)

	resetPasswordEntity, err := cmd.resetPasswordRepository.Find(resetpasswordrepo.FindResetPasswordDto{
		Ctx:       ctx,
		Token:     params.ResetPasswordToken,
		EagerLoad: []constants.EntityEagerLabel{constants.UserLabel},
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	actor := aggregates.NewUser(resetPasswordEntity.User)
	resetPassword := aggregates.NewResetPassword(resetPasswordEntity)

	if !resetPassword.Exists() {
		return cmd.messages.ResetPasswordNotExists(), nil
	}

	if !actor.Exists() {
		return cmd.messages.UserNotExists(), nil
	}

	if actor.IsActive() {
		return cmd.messages.UserIsAlreadyActive(), nil
	}

	if actor.IsBlocked() {
		return cmd.messages.UserIsBlocked(), nil
	}

	if resetPassword.IsExpired() {
		err := cmd.resetPasswordRepository.Delete(resetpasswordrepo.DeleteResetPasswordDto{
			Ctx:   ctx,
			Token: resetPasswordEntity.Token,
		})
		if err != nil {
			return cmd.messages.Default(), err
		}

		cmd.unitOfWork.Commit(ctx)

		return cmd.messages.InvalidResetToken(), nil
	}

	if !resetPassword.TokenIsValid(params.ResetPasswordHash) {
		return cmd.messages.InvalidResetToken(), nil
	}

	encrypted, err := hashing.Hash(params.NewPassword.Get())
	if err != nil {
		return cmd.messages.Default(), err
	}

	if err = cmd.userRepository.Update(userrepo.UpdateUserDto{
		Ctx:       ctx,
		UserToken: resetPasswordEntity.User.Token,
		Status:    &constants.UserActive,
		Password:  encrypted,
	}); err != nil {
		return cmd.messages.Default(), err
	}

	if err = cmd.resetPasswordRepository.Delete(resetpasswordrepo.DeleteResetPasswordDto{
		Ctx:   ctx,
		Token: resetPasswordEntity.Token,
	}); err != nil {
		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(ctx)

	return cmd.messages.Success(), nil
}
