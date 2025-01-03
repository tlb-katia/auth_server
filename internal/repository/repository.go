package repository

import (
	"context"
	"github.com/tlb_katia/auth/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=AuthRepository --output ./mocks --filename AuthRepository_mockery.go

type AuthRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
