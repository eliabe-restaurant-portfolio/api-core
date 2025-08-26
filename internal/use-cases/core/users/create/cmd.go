package createusercmd

import (
	"context"
	"time"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
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
	messages                    MessageProvider
	unitOfWork                  uow.UnitOfWork
	userRepository              userrepo.UserRepository
	resetPasswordRepository     resetpasswordrepo.ResetPasswordRepository
	sendInviteUserEmailProducer sendresetpasswordemailproducer.Producer
}

func New(adapters *adapters.Adapters) Command {
	return Command{
		messages:                    NewMessages(),
		unitOfWork:                  (*adapters).UnitOfWork(),
		userRepository:              (*adapters).Repositories().User(),
		resetPasswordRepository:     (*adapters).Repositories().ResetPassword(),
		sendInviteUserEmailProducer: (*adapters).Producers().SendPasswordResetEmail(),
	}
}

type Params struct {
	Context   context.Context
	Username  valueobjects.Username
	Email     valueobjects.Email
	TaxNumber valueobjects.TaxNumber
}

type Return struct {
	UserToken string
	Status    string
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	var err error
	var ctx = context.Background()
	cmd.unitOfWork.Init(ctx)

	existingEmailUser, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Email: &params.Email,
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	if existingEmailUser != nil {
		return cmd.messages.RepeatedUser(), nil
	}

	existingTaxNumberUser, err := cmd.userRepository.Find(userrepo.FindUserDto{
		TaxNumber: &params.TaxNumber,
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	if existingTaxNumberUser != nil {
		return cmd.messages.RepeatedUser(), nil
	}

	createdUser, err := cmd.userRepository.Create(userrepo.CreateUserDto{
		Ctx:       params.Context,
		Username:  params.Username,
		Email:     params.Email,
		TaxNumber: &params.TaxNumber,
		Status:    constants.UserInactive,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(params.Context)

		return cmd.messages.Default(), err
	}

	now := time.Now()
	future := now.Add(30 * time.Minute)

	encrypted, randomHash, err := hashing.GenerateRandom()
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

	if err = cmd.sendInviteUserEmailProducer.Send(sendresetpasswordemailproducer.SendPasswordResetEmailMessage{
		ResetPasswordToken: createdReset.Token.String(),
		RandomHash:         randomHash,
	}); err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(params.Context)

	return cmd.messages.Success(createdUser), nil
}
