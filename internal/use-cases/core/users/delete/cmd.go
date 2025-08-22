package deleteusercmd

import (
	"context"

	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	userrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/users"
	uow "github.com/eliabe-portfolio/restaurant-app/internal/unit-of-work"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/google/uuid"
)

type Command struct {
	messages       MessageProvider
	unitOfWork     uow.UnitOfWork
	userRepository userrepo.UserRepository
}

func New(repositories repositories.Provider, uow uow.UnitOfWork) Command {
	return Command{
		messages:       NewMessages(),
		unitOfWork:     uow,
		userRepository: repositories.User(),
	}
}

type Params struct {
	Context   context.Context
	UserToken uuid.UUID
}

type Return struct {
	UserToken string
	Status    string
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	var err error
	var ctx = context.Background()
	cmd.unitOfWork.Init(ctx)

	existing, err := cmd.userRepository.Find(userrepo.FindUserDto{
		Token: params.UserToken,
	})
	if err != nil {
		return cmd.messages.Default(), err
	}

	if existing == nil {
		return cmd.messages.UserNotExists(), nil
	}

	err = cmd.userRepository.Update(userrepo.UpdateUserDto{
		Ctx:       params.Context,
		UserToken: existing.Token,
	})
	if err != nil {
		cmd.unitOfWork.Rollback(ctx)

		return cmd.messages.Default(), err
	}

	cmd.unitOfWork.Commit(ctx)

	return cmd.messages.Success(), nil
}
