package resetpasswordcmd

import (
	"context"
	"time"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/aggregates"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	sendresetpasswordemailproducer "github.com/eliabe-restaurant-portfolio/api-core/internal/queues/producers/send-reset-password-email"
	resetpasswordrepo "github.com/eliabe-restaurant-portfolio/api-core/internal/repositories/reset-password"
	userrepo "github.com/eliabe-restaurant-portfolio/api-core/internal/repositories/users"
	uow "github.com/eliabe-restaurant-portfolio/api-core/internal/unit-of-work"
	valueobjects "github.com/eliabe-restaurant-portfolio/api-core/internal/value-objects"
	hashing "github.com/eliabe-restaurant-portfolio/api-core/pkg/hash"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
)

type Command struct {
	messages                       MessageProvider
	unitOfWork                     uow.UnitOfWork
	userRepository                 userrepo.UserRepository
	resetPasswordRepository        resetpasswordrepo.ResetPasswordRepository
	sendResetPasswordEmailProducer sendresetpasswordemailproducer.Producer
}

type Params struct {
	Email valueobjects.Email
}

func New(adapters *adapters.Adapters) Command {
	return Command{
		messages:                       NewMessages(),
		unitOfWork:                     (*adapters).UnitOfWork(),
		userRepository:                 (*adapters).Repositories().User(),
		resetPasswordRepository:        (*adapters).Repositories().ResetPassword(),
		sendResetPasswordEmailProducer: (*adapters).Producers().SendPasswordResetEmail(),
	}
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	var err error
	var ctx = context.Background()
	cmd.unitOfWork.Init(ctx)

	actorEntity, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Email:     &params.Email,
		EagerLoad: []constants.EntityEagerLabel{constants.ResetPasswordLabel},
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	var actor = aggregates.NewUser(actorEntity)

	if !actor.Exists() {
		return cmd.messages.UserNotExists(), nil
	}

	if len(actorEntity.ResetPasswords) != 0 {
		cmd.resetPasswordRepository.Delete(resetpasswordrepo.DeleteResetPasswordDto{
			Ctx:       ctx,
			UserToken: actorEntity.Token,
		})
	}

	encrypted, random, err := hashing.GenerateRandom()
	if err != nil {
		return cmd.messages.Default(), nil
	}

	created, err := cmd.resetPasswordRepository.Create(resetpasswordrepo.CreateResetPasswordDto{
		Ctx:       ctx,
		UserToken: actorEntity.Token,
		Hash:      encrypted,
		ValidAt:   time.Now().Add(30 * time.Minute),
	})
	if err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	err = cmd.userRepository.Update(userrepo.UpdateUserDto{
		Ctx:       ctx,
		UserToken: actorEntity.Token,
		Status:    &constants.UserInactive,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	err = cmd.sendResetPasswordEmailProducer.Send(sendresetpasswordemailproducer.SendPasswordResetEmailMessage{
		ResetPasswordToken: created.Token.String(),
		RandomHash:         random,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(ctx)

	return cmd.messages.Success(), nil
}
