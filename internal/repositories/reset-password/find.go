package resetpasswordrepo

import (
	"context"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FindResetPasswordDto struct {
	Ctx       context.Context
	Token     uuid.UUID
	EagerLoad []constants.EntityEagerLabel
}

func (r repository) Find(dto FindResetPasswordDto) (*entities.ResetPassword, error) {
	var zero uuid.UUID
	var user *entities.ResetPassword
	var query = r.conn

	if dto.Token != zero {
		if err := query.Where("token = ?", dto.Token).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, errs.New(err.Error())
		}
	}

	for _, name := range dto.EagerLoad {
		query = query.Preload(string(name))
	}

	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errs.New(err.Error())
	}

	return user, nil
}
