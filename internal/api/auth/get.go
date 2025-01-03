package auth

import (
	"context"
	"github.com/tlb_katia/auth/internal/converter"
	"github.com/tlb_katia/auth/pkg/auth_v1"
)

func (i *Implementation) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	user, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &auth_v1.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
