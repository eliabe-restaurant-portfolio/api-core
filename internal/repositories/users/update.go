package userrepo

import (
	"context"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/entities"
	"github.com/eliabe-restaurant-portfolio/api-core/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateUserDto struct {
	Ctx                 context.Context
	UserToken           uuid.UUID
	FailedLoginAttempts *int
	Status              *constants.UserStatus
	Password            string
}

func (r repository) Update(dto UpdateUserDto) error {
	tx := r.conn
	if transaction, ok := dto.Ctx.Value(constants.TransactionKey).(*gorm.DB); ok {
		tx = transaction
	}

	updates := make(map[string]interface{})

	if dto.FailedLoginAttempts != nil {
		updates["failed_login_attempts"] = *dto.FailedLoginAttempts
	}
	if dto.Status != nil {
		updates["status"] = *dto.Status
	}
	if dto.Password != "" {
		updates["password"] = dto.Password
	}

	if len(updates) == 0 {
		return errs.New("it has not updates")
	}

	result := tx.WithContext(dto.Ctx).
		Model(&entities.User{}).
		Where("token = ?", dto.UserToken).
		Updates(updates)

	if result.Error != nil {
		return errs.New(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
