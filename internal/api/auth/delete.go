package auth

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tlb_katia/auth/pkg/auth_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*empty.Empty, error) {
	err := i.authService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
