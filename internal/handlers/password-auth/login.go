package authhdl

import (
	"log"
	"net/http"

	traditionallogincmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/auth/password-auth/login"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
)

type LoginHttpDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (hdl AuthHandler) Login(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto LoginHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	params, err := parseLoginParams(dto)
	if err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}

	result, err := traditionallogincmd.New(hdl.adapters).Execute(*params)

	if err != nil {
		log.Printf("password login command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func parseLoginParams(dto LoginHttpDto) (*traditionallogincmd.Params, error) {
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
