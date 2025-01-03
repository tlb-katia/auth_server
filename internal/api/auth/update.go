package auth

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tlb_katia/auth/internal/converter"
	"github.com/tlb_katia/auth/pkg/auth_v1"
)

func (i *Implementation) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*empty.Empty, error) {
	err := i.authService.Update(ctx, converter.ToUserUpdateFromGRPC(req.GetUser()))
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
