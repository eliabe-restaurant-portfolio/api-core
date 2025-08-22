package authhdl

import (
	"log"
	"net/http"

	changepasswordcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/change-password"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/request"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ChangePasswordHttpDto struct {
	Email       string `json:"email" binding:"required,email"`
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
	Token       string `json:"token" binding:"required"`
}

func (dto *ChangePasswordHttpDto) Validate() error {
	var validate = validator.New()
	return validate.Struct(dto)
}

func (hdl AuthHandler) ChangePassword(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ChangePasswordHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	if err := dto.Validate(); err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	actorId, err := request.RecoveryActor(ctx)
	if err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := buildChangePasswordParams(actorId, dto)
	if err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := changepasswordcmd.New(
		hdl.repositories,
		hdl.uow,
	).Execute(*params)
	if err != nil {
		log.Printf("change password command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func buildChangePasswordParams(actorId uuid.UUID, dto ChangePasswordHttpDto) (*changepasswordcmd.Params, error) {
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

	uuid, err := uuid.Parse(dto.Token)
	if err != nil {
		return nil, err
	}

	return &changepasswordcmd.Params{
		ActorToken:  uuid,
		Email:       email,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}, nil
}
