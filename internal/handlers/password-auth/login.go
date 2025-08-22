package authhdl

import (
	"log"
	"net/http"

	traditionallogincmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/login"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoginHttpDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (dto *LoginHttpDto) Validate() error {
	var validate = validator.New()
	return validate.Struct(dto)
}

func (hdl AuthHandler) Login(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto LoginHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	if err := dto.Validate(); err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := buildLoginParams(dto)
	if err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := traditionallogincmd.New(
		hdl.repositories,
		hdl.uow,
	).Execute(*params)
	if err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func buildLoginParams(dto LoginHttpDto) (*traditionallogincmd.Params, error) {
	email, err := valueobjects.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	password, err := valueobjects.NewPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	return &traditionallogincmd.Params{
		Email:    email,
		Password: password,
	}, nil
}
