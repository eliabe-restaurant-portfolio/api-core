package aggregates

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
)

type User interface {
	Exists() bool
	IsActive() bool
	IsBlocked() bool
	IsInactive() bool
	HasMaxFailedLoginAttempts() bool
	HasNotFailedLoginAttemps() bool
	GetFailedLoginAttempts() int
	IsValidPassword(password valueobjects.Password) bool
}

type userAggr struct {
	val *entities.User
}

func NewUser(user *entities.User) User {
	return userAggr{val: user}
}

func (aggr userAggr) Exists() bool {
	return aggr.val != nil
}

func (aggr userAggr) IsBlocked() bool {
	return aggr.val.Status == constants.UserBloqued
}

func (aggr userAggr) IsActive() bool {
	return aggr.val.Status == constants.UserActive
}

func (aggr userAggr) IsInactive() bool {
	return aggr.val.Status == constants.UserInactive
}

func (aggr userAggr) HasMaxFailedLoginAttempts() bool {
	return aggr.val.FailedLoginAttempts > constants.MAX_FAILED_LOGIN_ATTEMPTS
}

func (aggr userAggr) HasNotFailedLoginAttemps() bool {
	return aggr.val.FailedLoginAttempts == 0
}

func (aggr userAggr) GetFailedLoginAttempts() int {
	return aggr.val.FailedLoginAttempts
}

func (aggr userAggr) IsValidPassword(password valueobjects.Password) bool {
	return aggr.val.Password == password.Get()
}
