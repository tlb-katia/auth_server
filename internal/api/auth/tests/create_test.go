package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/tlb_katia/auth/internal/api/auth"
	"github.com/tlb_katia/auth/internal/model"
	"github.com/tlb_katia/auth/internal/service"
	"github.com/tlb_katia/auth/internal/service/mocks"
	"github.com/tlb_katia/auth/pkg/auth_v1"
	"testing"
)

func randomRole() auth_v1.Role {
	return auth_v1.Role(gofakeit.Number(0, 1))
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *auth_v1.CreateRequest
	}

	type authServiceMockFunc func() service.AuthService

	var (
		ctx = context.Background()

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(false, false, false, false, false, 7)
		role     = randomRole()

		serviceErr = fmt.Errorf("service error")

		req = &auth_v1.CreateRequest{
			User: &auth_v1.UserInfo{
				Name:     name,
				Email:    email,
				Password: password,
				Role:     role,
			},
		}

		reqInfo = &model.UserInfo{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     model.Role(role),
		}

		expected = &auth_v1.CreateResponse{
			Id: id,
		}
	)
	tests := []struct {
		name            string
		args            args
		want            *auth_v1.CreateResponse
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "Successful Case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: expected,
			err:  nil,
			authServiceMock: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Create", ctx, reqInfo).Return(id, nil)

				return mock
			},
		},
		{
			name: "Failed Case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Create", ctx, reqInfo).Return(int64(0), serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			authServiceMock := tt.authServiceMock()
			api := auth.NewImplementation(authServiceMock)

			newId, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newId)
		})
	}

}
