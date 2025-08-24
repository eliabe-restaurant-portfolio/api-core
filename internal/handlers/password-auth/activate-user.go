package authhdl

import (
	"log"
	"net/http"

	activateusercmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/activate-user"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivateLoginHttpDto struct {
	ResetPasswordToken string `json:"reset_password_token"`
	ResetPasswordHash  string `json:"reset_password_hash"`
	NewPassword        string `json:"new_password"`
}

func (hdl AuthHandler) ActivateUser(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ActivateLoginHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := parseActivateUser(dto)
	if err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := activateusercmd.New(hdl.adapters).Execute(*params)

	if err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func parseActivateUser(dto ActivateLoginHttpDto) (*activateusercmd.Params, error) {
	newPassword, err := valueobjects.NewPassword(dto.NewPassword)
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.Parse(dto.ResetPasswordToken)
	if err != nil {
		return nil, err
	}

	return &activateusercmd.Params{
		ResetPasswordToken: uuid,
		ResetPasswordHash:  dto.ResetPasswordToken,
		NewPassword:        newPassword,
	}, nil
}
