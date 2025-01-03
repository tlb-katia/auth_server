package auth

import (
	"context"
	"github.com/tlb_katia/auth/internal/converter"
	"github.com/tlb_katia/auth/pkg/auth_v1"
)

func (i *Implementation) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	id, err := i.authService.Create(ctx, converter.ToUserInfoFromGRPC(req.GetUser()))

	if err != nil {
		return nil, err
	}
	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}
