package sendresetpasswordemailcmd

import (
	"fmt"

	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	resetpasswordrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/reset-password"
	"github.com/eliabe-portfolio/restaurant-app/pkg/email"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/google/uuid"
)

type Command struct {
	messageProvider         MessageProvider
	resetPasswordRepository resetpasswordrepo.ResetPasswordRepository
}

func New(adapters *adapters.Adapters) Command {
	return Command{
		messageProvider:         messageProvider,
		resetPasswordRepository: (*adapters).Repositories().ResetPassword(),
	}
}

type Params struct {
	ResetPasswordToken uuid.UUID
	Token              string
}

type RelatedEntities struct {
	PasswordResetToken *entities.ResetPassword
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	related, err := cmd.getRelatedEntities(params)
	if err != nil {
		return returns.Api{}, err
	}

	if related.PasswordResetToken == nil {
		return returns.Api{}, fmt.Errorf("password reset token not found for uuid: %s", params.ResetPasswordToken)
	}

	user := related.PasswordResetToken.User

	err = email.New().SendPasswordResetEmail(email.PasswordResetEmailInput{
		To:       user.Email,
		Subject:  "Redefinição de Senha",
		UserName: user.Username,
		Url:      params.Token,
	})
	if err != nil {
		return returns.Api{}, err
	}

	return cmd.messageProvider.Success(), nil
}

func (cmd Command) getRelatedEntities(params Params) (*RelatedEntities, error) {
	related := &RelatedEntities{}

	passwordResetToken, err := cmd.resetPasswordRepository.Find(resetpasswordrepo.FindResetPasswordDto{
		Token:     params.ResetPasswordToken,
		EagerLoad: []constants.EntityEagerLabel{constants.UserLabel},
	})
	if err != nil {
		return nil, err
	}

	related.PasswordResetToken = passwordResetToken

	return related, nil
}
