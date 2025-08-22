package authhdl

import (
	"log"
	"net/http"

	activatelogincmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/activate-login"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ActivateLoginHttpDto struct {
	ResetPasswordToken string `json:"reset_password_token" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
}

func (dto *ActivateLoginHttpDto) Validate() error {
	var validate = validator.New()
	return validate.Struct(dto)
}

func (hdl AuthHandler) ActivateLogin(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ActivateLoginHttpDto

	hash := ctx.Query("hash")

	log.Println(hash)

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	if err := dto.Validate(); err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := builActivateLoginParams(dto, hash)
	if err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := activatelogincmd.New(
		hdl.repositories,
		hdl.uow,
	).Execute(*params)

	if err != nil {
		log.Printf("activate login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func builActivateLoginParams(dto ActivateLoginHttpDto, hash string) (*activatelogincmd.Params, error) {
	newPassword, err := valueobjects.NewPassword(dto.NewPassword)
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.Parse(dto.ResetPasswordToken)
	if err != nil {
		return nil, err
	}

	return &activatelogincmd.Params{
		ResetPasswordToken: uuid,
		NewPassword:        newPassword,
		ResetPasswordHash:  hash,
	}, nil
}
