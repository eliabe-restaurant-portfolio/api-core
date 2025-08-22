package aggregates

import (
	"time"

	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	hashing "github.com/eliabe-portfolio/restaurant-app/pkg/hash"
)

type ResetPassword interface {
	Exists() bool
	IsExpired() bool
	TokenIsValid(input string) bool
}

type resetPasswordAggr struct {
	val *entities.ResetPassword
}

func NewResetPassword(reset *entities.ResetPassword) ResetPassword {
	return resetPasswordAggr{val: reset}
}

func (aggr resetPasswordAggr) Exists() bool {
	return aggr.val != nil
}

func (aggr resetPasswordAggr) IsExpired() bool {
	return time.Now().Before(aggr.val.CreatedAt)
}

func (aggr resetPasswordAggr) TokenIsValid(input string) bool {
	return hashing.Compare(aggr.val.Hash, input)
}
