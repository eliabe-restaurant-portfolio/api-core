package userrepo

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	valueobjects "github.com/eliabe-restaurant-portfolio/api-core/internal/value-objects"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FindUserDto struct {
	Token     uuid.UUID
	Email     *valueobjects.Email
	Username  *valueobjects.Username
	TaxNumber *valueobjects.TaxNumber
	EagerLoad []constants.EntityEagerLabel
}

func (r repository) Find(dto FindUserDto) (*entities.User, error) {
	var zero uuid.UUID
	var user *entities.User
	var query = r.conn

	if dto.Token != zero {
		if err := query.Where("token = ?", dto.Token.String()).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, errs.New(err.Error())
		}
	}

	if dto.Email != nil {
		if err := query.Where("email = ?", dto.Email.Get()).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, errs.New(err.Error())
		}
	}

	if dto.Username != nil {
		if err := query.Where("username = ?", dto.Username.Get()).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, errs.New(err.Error())
		}
	}

	if dto.TaxNumber != nil {
		if err := query.Where("tax_number = ?", dto.Username.Get()).First(&user).Error; err != nil {
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
