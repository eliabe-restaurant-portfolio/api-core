package userhdl

import (
	"context"
	"log"
	"net/http"

	deleteusercmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/core/users/delete"
	"github.com/eliabe-portfolio/restaurant-app/pkg/request"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
)

func (hdl UserHandler) Delete(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})

	actor, err := request.RecoveryActor(ctx)
	if err != nil {
		log.Printf("delete user command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	params := deleteusercmd.Params{
		Context:   context.Background(),
		UserToken: actor,
	}

	result, err := deleteusercmd.New(
		hdl.repositories,
		hdl.uow,
	).Execute(params)
	if err != nil {
		log.Printf("delete user command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	}

	log.Printf("delete user command: %v", result)

	ctx.JSON(http.StatusOK, result)
}
