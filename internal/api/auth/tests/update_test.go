package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/require"
	"github.com/tlb_katia/auth/internal/api/auth"
	"github.com/tlb_katia/auth/internal/model"
	"github.com/tlb_katia/auth/internal/service"
	"github.com/tlb_katia/auth/internal/service/mocks"
	"github.com/tlb_katia/auth/pkg/auth_v1"
	"testing"
)

func TestUpdate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *auth_v1.UpdateRequest
	}

	type serviceMockFunc func() service.AuthService

	var (
		ctx = context.Background()

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		req = &auth_v1.UpdateRequest{
			User: &auth_v1.UserUpdate{
				Id:    id,
				Name:  &wrappers.StringValue{Value: name},
				Email: &wrappers.StringValue{Value: email},
			},
		}

		reqServ = &model.UserUpdate{
			ID:    id,
			Name:  &wrappers.StringValue{Value: name},
			Email: &wrappers.StringValue{Value: email},
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		serviceMockFunc serviceMockFunc
	}{
		{
			name: "Success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &empty.Empty{},
			err:  nil,
			serviceMockFunc: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Update", ctx, reqServ).Return(nil)

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
			err:  serviceErr,
			serviceMockFunc: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Update", ctx, reqServ).Return(serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.serviceMockFunc()
			api := auth.NewImplementation(authServiceMock)
			apiRes, apiErr := api.Update(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, apiErr)
		})
	}
}
