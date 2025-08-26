package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResetPassword struct {
	Token     uuid.UUID      `json:"token" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserToken uuid.UUID      `json:"user_token" gorm:"type:uuid;index"`
	Hash      string         `json:"hash" gorm:"type:string"`
	ValidAt   time.Time      `json:"valid_at"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User *User `gorm:"foreignKey:UserToken;references:Token"`
}
