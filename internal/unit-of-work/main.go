package uow

import (
	"context"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/connections"
	"gorm.io/gorm"
)

type UnitOfWork interface {
	Init(ctx context.Context)
	Commit(ctx context.Context)
	Rollback(ctx context.Context)
}

type uow struct {
	connection *gorm.DB
	tx         *gorm.DB
}

func New(conns *connections.Provider) UnitOfWork {
	return &uow{
		connection: conns.GetPostgres(),
	}
}

func (u *uow) Init(ctx context.Context) {
	tx := u.connection.WithContext(ctx).Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	u.tx = tx
}

func (u *uow) Commit(ctx context.Context) {
	tx := u.connection.WithContext(ctx).Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	if err := u.tx.Commit().Error; err != nil {
		panic(err)
	}
}

func (u *uow) Rollback(ctx context.Context) {
	tx := u.connection.WithContext(ctx).Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	if err := u.tx.Rollback().Error; err != nil {
		panic(err)
	}
}
