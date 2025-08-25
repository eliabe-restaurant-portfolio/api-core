package resetpasswordrepo

import (
	"context"
	"time"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateResetPasswordDto struct {
	Ctx       context.Context
	UserToken uuid.UUID
	Hash      string
	ValidAt   time.Time
}

func (r repository) Create(dto CreateResetPasswordDto) (*entities.ResetPassword, error) {
	tx := r.conn
	if transaction, ok := dto.Ctx.Value(constants.TransactionKey).(*gorm.DB); ok {
		tx = transaction
	}

	item := entities.ResetPassword{
		UserToken: dto.UserToken,
		Hash:      dto.Hash,
		ValidAt:   dto.ValidAt,
	}

	if err := tx.Create(&item).Error; err != nil {
		return nil, errs.New(err.Error())
	}

	return &item, nil
}
