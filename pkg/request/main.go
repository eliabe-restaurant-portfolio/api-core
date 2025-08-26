package request

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RecoveryActor(ctx *gin.Context) (uuid.UUID, error) {
	var zero uuid.UUID
	actor, exists := ctx.Get("actor")
	if !exists {
		return zero, fmt.Errorf("actor not found in context")
	}

	v := actor.(string)

	uuid, err := uuid.Parse(v)
	if err != nil {
		return zero, fmt.Errorf("error to transform uuid")
	}

	return uuid, nil
}
