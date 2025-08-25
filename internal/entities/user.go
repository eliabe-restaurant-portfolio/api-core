package entities

import (
	"time"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Token               uuid.UUID            `json:"token" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username            string               `json:"username" gorm:"not null;uniqueIndex"`
	Email               string               `json:"email" gorm:"not null;uniqueIndex"`
	TaxNumber           string               `json:"tax_number" gorm:"not null;uniqueIndex"`
	Password            string               `json:"password" gorm:"not null"`
	Status              constants.UserStatus `json:"status" gorm:"not null"`
	FailedLoginAttempts int                  `json:"failed_login_attempts" gorm:"not null;default:0"`
	CreatedAt           time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt       `json:"-" gorm:"index"` // Usar gorm.DeletedAt

	ResetPasswords []ResetPassword `gorm:"foreignKey:UserToken;references:Token"`
}
