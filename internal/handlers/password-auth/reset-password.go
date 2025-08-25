package authhdl

import (
	"log"
	"net/http"

	resetpasswordcmd "github.com/eliabe-restaurant-portfolio/api-core/internal/use-cases/auth/password-auth/reset-password"
	valueobjects "github.com/eliabe-restaurant-portfolio/api-core/internal/value-objects"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/returns"
	"github.com/gin-gonic/gin"
)

type ResetPasswordHttpDto struct {
	Email string `json:"email"`
}

func (hdl AuthHandler) RequestResetPassword(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ResetPasswordHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := parseRequestResetPassword(dto)
	if err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := resetpasswordcmd.New(hdl.adapters).Execute(*params)

	if err != nil {
		log.Printf("reset password command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	log.Printf("token de redefinição de senha enviado com sucesso para o email %s", dto.Email)

	ctx.JSON(http.StatusOK, result)
}

func parseRequestResetPassword(dto ResetPasswordHttpDto) (*resetpasswordcmd.Params, error) {
	email, err := valueobjects.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	return &resetpasswordcmd.Params{
		Email: email,
	}, nil
}
