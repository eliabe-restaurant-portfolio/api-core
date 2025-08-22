package authhdl

import (
	"log"
	"net/http"

	resetpasswordcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/reset-password"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ResetPasswordHttpDto struct {
	Email string `json:"email" binding:"required,email"`
}

func (dto *ResetPasswordHttpDto) Validate() error {
	var validate = validator.New()
	return validate.Struct(dto)
}

func (hdl AuthHandler) ResetPassword(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ResetPasswordHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	if err := dto.Validate(); err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := buildResetPasswordParams(dto)
	if err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := resetpasswordcmd.New(
		hdl.repositories,
		hdl.uow,
		hdl.producers,
	).Execute(*params)
	if err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	log.Printf("token de redefinição de senha enviado com sucesso para o email %s", dto.Email)

	ctx.JSON(http.StatusOK, result)
}

func buildResetPasswordParams(dto ResetPasswordHttpDto) (*resetpasswordcmd.Params, error) {
	email, err := valueobjects.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	return &resetpasswordcmd.Params{
		Email: email,
	}, nil
}
