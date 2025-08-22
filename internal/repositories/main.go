package repositories

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	resetpasswordrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/reset-password"
	userrepo "github.com/eliabe-portfolio/restaurant-app/internal/repositories/users"
)

type Provider interface {
	User() userrepo.UserRepository
	ResetPassword() resetpasswordrepo.ResetPasswordRepository
}

type repositories struct {
	user          userrepo.UserRepository
	resetPassword resetpasswordrepo.ResetPasswordRepository
}

func New(connections *connections.Provider) Provider {
	return repositories{
		user:          userrepo.New(connections.Postgres.Get()),
		resetPassword: resetpasswordrepo.New(connections.GetPostgres()),
	}
}

func (r repositories) User() userrepo.UserRepository {
	return r.user
}

func (r repositories) ResetPassword() resetpasswordrepo.ResetPasswordRepository {
	return r.resetPassword
}
