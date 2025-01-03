package service

import (
	"context"
	"github.com/tlb_katia/auth/internal/model"
)

type AuthService interface {
	Create(ctx context.Context, user *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
