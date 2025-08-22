package resetpasswordrepo

import (
	"context"

	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/eliabe-portfolio/restaurant-app/internal/entities"
	"github.com/eliabe-portfolio/restaurant-app/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteResetPasswordDto struct {
	Ctx       context.Context
	Token     uuid.UUID
	UserToken uuid.UUID
}

func (r repository) Delete(dto DeleteResetPasswordDto) error {
	var zero uuid.UUID
	tx := r.conn
	if transaction, ok := dto.Ctx.Value(constants.TransactionKey).(*gorm.DB); ok {
		tx = transaction
	}

	if dto.Token == zero && dto.UserToken == zero {
		return errs.New("it has not filters to delete")
	}

	query := tx.Model(&entities.ResetPassword{})

	if dto.Token != zero {
		query.Where("token = ?", dto.Token)
	}

	if dto.UserToken != zero {
		query.Where("user_token = ?", dto.Token)
	}

	if err := query.Delete(&entities.ResetPassword{}).Error; err != nil {
		return errs.New(err.Error())
	}

	return nil
}
