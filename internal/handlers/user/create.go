package userhdl

import (
	"context"
	"log"
	"net/http"

	createusercmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/core/users/create"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateUserHttpDto struct {
	Email     string  `json:"email" validate:"required,email"`
	Username  string  `json:"username" validate:"required,min=6"`
	TaxNumber *string `json:"tax_number"`
}

func (dto *CreateUserHttpDto) Validate() error {
	var validate = validator.New()
	return validate.Struct(dto)
}

func (hdl UserHandler) Create(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto CreateUserHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("create user command: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	if err := dto.Validate(); err != nil {
		log.Printf("create user command: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	params, err := buildCreateUserParams(dto)
	if err != nil {
		log.Printf("create user command: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	result, err := createusercmd.New(hdl.adapters).Execute(*params)
	if err != nil {
		log.Printf("create user command: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	log.Printf("create user command: %v", result)
	ctx.JSON(result.Code, result)
}

func buildCreateUserParams(dto CreateUserHttpDto) (*createusercmd.Params, error) {
	var taxNumber valueobjects.TaxNumber
	email, err := valueobjects.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	if dto.TaxNumber != nil {
		taxNumber, err = valueobjects.NewTaxNumber(*dto.TaxNumber)
		if err != nil {
			return nil, err
		}
	}

	username, err := valueobjects.NewUsername(dto.Username)
	if err != nil {
		return nil, err
	}

	return &createusercmd.Params{
		Context:   context.Background(),
		Email:     email,
		TaxNumber: &taxNumber,
		Username:  username,
	}, nil
}
