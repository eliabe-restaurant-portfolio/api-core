package userrepo

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(dto CreateUserDto) (*entities.User, error)
	Find(dto FindUserDto) (*entities.User, error)
	Update(dto UpdateUserDto) error
}

type repository struct {
	conn  *gorm.DB
	table string
}

func New(conn *gorm.DB) UserRepository {
	var table = "users"
	return repository{conn: conn, table: table}
}
