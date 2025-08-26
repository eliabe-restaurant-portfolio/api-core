package resetpasswordrepo

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"gorm.io/gorm"
)

type ResetPasswordRepository interface {
	Find(dto FindResetPasswordDto) (*entities.ResetPassword, error)
	Delete(dto DeleteResetPasswordDto) error
	Create(dto CreateResetPasswordDto) (*entities.ResetPassword, error)
}

type repository struct {
	conn  *gorm.DB
	table string
}

func New(conn *gorm.DB) ResetPasswordRepository {
	var table = "reset_passwords"
	return repository{conn: conn, table: table}
}
