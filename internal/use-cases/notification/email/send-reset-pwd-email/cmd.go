package sendresetpasswordemailcmd

import (
	"fmt"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	resetpasswordrepo "github.com/eliabe-restaurant-portfolio/api-core/internal/repositories/reset-password"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/email"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
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
	RandomHash         string
	RandomPassword     string
}

func (cmd Command) Execute(params Params) (returns.Api, error) {
	resetPassword, err := cmd.resetPasswordRepository.Find(resetpasswordrepo.FindResetPasswordDto{
		Token:     params.ResetPasswordToken,
		EagerLoad: []constants.EntityEagerLabel{constants.UserLabel},
	})
	if err != nil {
		return cmd.messageProvider.Default(), err
	}

	if resetPassword == nil {
		return returns.Api{}, fmt.Errorf("password reset token not found for uuid: %s", params.ResetPasswordToken)
	}

	err = email.New().SendPasswordResetEmail(email.PasswordResetEmailInput{
		To:                 resetPassword.User.Email,
		Subject:            "Redefinição de Senha",
		UserName:           resetPassword.User.Username,
		Hash:               params.RandomHash,
		ResetPasswordToken: resetPassword.Token.String(),
	})
	if err != nil {
		return cmd.messageProvider.Default(), err
	}

	return cmd.messageProvider.Success(), nil
}
