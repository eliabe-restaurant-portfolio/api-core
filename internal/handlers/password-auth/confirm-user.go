package authhdl

import (
	"log"
	"net/http"

	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/gin-gonic/gin"
)

type ConfirmUserHttpDto struct {
	UserConfirmCode string `json:"user_confirm_code"`
	UserToken       string `json:"user_token"`
}

func (hdl AuthHandler) ConfirmUser(ctx *gin.Context) {
	var defaultError = returns.InternalServerError([]string{})
	var dto ChangePasswordHttpDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("confirm user command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	}
	/* params, err := parseChangePassword(actorToken, dto)
	if err != nil {
		log.Printf("confirm user command: %v", err)
		ctx.JSON(http.StatusBadRequest, defaultError)
		return
	} */

	/* result, err := changepasswordcmd.New(hdl.adapters).Execute(*params)
	if err != nil {
		log.Printf("confirm user command: %v", err)
		ctx.JSON(http.StatusInternalServerError, defaultError)
		return
	} */

	/* ctx.JSON(http.StatusOK, result) */
}
