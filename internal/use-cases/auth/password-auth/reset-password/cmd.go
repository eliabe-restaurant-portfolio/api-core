package resetpasswordcmd

import (
	"context"
	"time"

	"github.com/eliabe-portfolio/restaurant-app/internal/aggregates"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	sendresetpasswordemailconsumer "github.com/eliabe-portfolio/restaurant-app/internal/queues/consumers/send-reset-password-email"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/producers"
	sendresetpasswordemailproducer "github.com/eliabe-portfolio/restaurant-app/internal/queues/producers/send-reset-password-email"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	resetpasswordrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/reset-password"
	userrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/users"
	uow "github.com/eliabe-portfolio/restaurant-app/internal/unit-of-work"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	hashing "github.com/eliabe-portfolio/restaurant-app/pkg/hash"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
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

func New(repositories repositories.Provider, uow uow.UnitOfWork, producers producers.Provider) Command {
	return Command{
		messages:                       NewMessages(),
		unitOfWork:                     uow,
		userRepository:                 repositories.User(),
		resetPasswordRepository:        repositories.ResetPassword(),
		sendResetPasswordEmailProducer: producers.SendPasswordResetEmail(),
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
		ValidAt:   time.Now(),
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

	err = cmd.sendResetPasswordEmailProducer.Send(sendresetpasswordemailconsumer.SendPasswordResetEmailMessage{
		ResetPasswordToken: created.Token.String(),
		Token:              random,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(ctx)

	return cmd.messages.Success(), nil
}
