package userrepo

import (
	"context"

	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	valueobjects "github.com/eliabe-portfolio/restaurant-app/internal/value-objects"
	"github.com/eliabe-portfolio/restaurant-app/pkg/errs"
	"gorm.io/gorm"
)

type CreateUserDto struct {
	Ctx       context.Context
	Username  valueobjects.Username
	Email     valueobjects.Email
	Password  string
	TaxNumber *valueobjects.TaxNumber
	Status    constants.UserStatus
}

func (r repository) Create(dto CreateUserDto) (*entities.User, error) {
	var tx = r.conn
	var taxNumber *string = nil
	if transaction, ok := dto.Ctx.Value(constants.TransactionKey).(*gorm.DB); ok {
		tx = transaction
	}

	if dto.TaxNumber != nil {
		val := dto.TaxNumber.Get()
		taxNumber = &val
	}

	item := entities.User{
		Email:     dto.Email.Get(),
		Username:  dto.Username.Get(),
		Password:  dto.Password,
		Status:    dto.Status,
		TaxNumber: *taxNumber,
	}

	if err := tx.Create(&item).Error; err != nil {
		return nil, errs.New(err.Error())
	}

	return &item, nil
}
