package authhdl

import (
	"log"
	"net/http"

	changepasswordcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/change-password"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/request"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChangePasswordHttpDto struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (hdl AuthHandler) ChangePassword(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ChangePasswordHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	actorToken, err := request.RecoveryActor(ctx)
	if err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := parseChangePassword(actorToken, dto)
	if err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := changepasswordcmd.New(hdl.adapters).Execute(*params)
	if err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func parseChangePassword(actorToken uuid.UUID, dto ChangePasswordHttpDto) (*changepasswordcmd.Params, error) {
	email, err := valueobjects.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	oldPassword, err := valueobjects.NewPassword(dto.OldPassword)
	if err != nil {
		return nil, err
	}

	newPassword, err := valueobjects.NewPassword(dto.NewPassword)
	if err != nil {
		return nil, err
	}

	return &changepasswordcmd.Params{
		ActorToken:  actorToken,
		Email:       email,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}, nil
}
