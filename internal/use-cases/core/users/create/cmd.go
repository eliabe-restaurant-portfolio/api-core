package createusercmd

import (
	"context"
	"time"

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
	messages                    MessageProvider
	unitOfWork                  uow.UnitOfWork
	userRepository              userrepo.UserRepository
	resetPasswordRepository     resetpasswordrepo.ResetPasswordRepository
	sendInviteUserEmailProducer sendresetpasswordemailproducer.Producer
}

func New(repositories repositories.Provider, uow uow.UnitOfWork, producers producers.Provider) Command {
	return Command{
		messages:                    NewMessages(),
		unitOfWork:                  uow,
		userRepository:              repositories.User(),
		resetPasswordRepository:     repositories.ResetPassword(),
		sendInviteUserEmailProducer: producers.SendPasswordResetEmail(),
	}
}

type Params struct {
	Context   context.Context
	Username  valueobjects.Username
	Email     valueobjects.Email
	TaxNumber *valueobjects.TaxNumber
}

type Return struct {
	UserToken string
	Status    string
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	var err error
	var ctx = context.Background()
	cmd.unitOfWork.Init(ctx)

	existingEmail, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Email: &params.Email,
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	if existingEmail != nil {
		return cmd.messages.ExistsUserWithSameEmail(), nil
	}

	existingUsername, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Email: &params.Email,
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	if existingUsername != nil {
		return cmd.messages.ExistsUserWithSameEmail(), nil
	}

	// random, err := valueobjects.NewRandonPassword()
	// if err != nil {
	// 	return cmd.messages.Default(), err
	// }

	pass, err := valueobjects.NewPassword("#Eli2025")
	if err != nil {
		return cmd.messages.Default(), err
	}

	createdUser, err := cmd.userRepository.Create(userrepo.CreateUserDto{
		Ctx:       params.Context,
		Username:  params.Username,
		Email:     params.Email,
		Password:  pass,
		TaxNumber: params.TaxNumber,
		Status:    constants.UserInactive,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(params.Context)

		return cmd.messages.Default(), err
	}

	now := time.Now()
	future := now.Add(30 * time.Minute)

	encrypted, random, err := hashing.GenerateRandom()
	if err != nil {
		cmd.unitOfWork.Rollback(params.Context)

		return cmd.messages.Default(), err
	}

	createdReset, err := cmd.resetPasswordRepository.Create(resetpasswordrepo.CreateResetPasswordDto{
		Ctx:       params.Context,
		UserToken: createdUser.Token,
		Hash:      encrypted,
		ValidAt:   future,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(params.Context)

		return cmd.messages.Default(), err
	}

	if err = cmd.sendInviteUserEmailProducer.Send(sendresetpasswordemailconsumer.SendPasswordResetEmailMessage{
		ResetPasswordToken: createdReset.Token.String(),
		Token:              random,
	}); err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(params.Context)

	return cmd.messages.Success(createdUser), nil
}
