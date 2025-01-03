package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"github.com/tlb_katia/auth/internal/api/auth"
	"github.com/tlb_katia/auth/internal/service"
	"github.com/tlb_katia/auth/internal/service/mocks"
	"github.com/tlb_katia/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestDelete(t *testing.T) {
	type serviceMockFunc func() service.AuthService
	type args struct {
		ctx context.Context
		req *auth_v1.DeleteRequest
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		req = &auth_v1.DeleteRequest{
			Id: id,
		}

		servReq = id

		servError = fmt.Errorf("server error")
	)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		serviceMockFunc serviceMockFunc
	}{
		{
			name: "Successful case",
			args: args{
				req: req,
				ctx: ctx,
			},
			want: &emptypb.Empty{},
			err:  nil,
			serviceMockFunc: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Delete", ctx, servReq).Return(nil)

				return mock
			},
		},
		{
			name: "Failed case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  servError,
			serviceMockFunc: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Delete", ctx, servReq).Return(servError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			serviceMock := tt.serviceMockFunc()
			api := auth.NewImplementation(serviceMock)

			apiRes, apiErr := api.Delete(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, apiErr)
		})
	}
}
