package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/require"
	"github.com/tlb_katia/auth/internal/api/auth"
	"github.com/tlb_katia/auth/internal/model"
	"github.com/tlb_katia/auth/internal/service"
	"github.com/tlb_katia/auth/internal/service/mocks"
	"github.com/tlb_katia/auth/pkg/auth_v1"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
		req *auth_v1.GetRequest
	}

	type authServerMocFunc func() service.AuthService

	var (
		ctx = context.Background()

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = randomRole()

		currentTime = time.Now()
		ts          = &timestamp.Timestamp{
			Seconds: currentTime.Unix(),              // Количество секунд с 1970-01-01
			Nanos:   int32(currentTime.Nanosecond()), // Количество наносекунд
		}
		createdAt = ts
		updatedAt = ts

		serviceErr = fmt.Errorf("service error")

		req = &auth_v1.GetRequest{
			Id: id,
		}

		servReq  = id
		servResp = &model.User{
			ID: id,
			User: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  model.Role(role),
			},
			CreatedAt: currentTime,
			UpdatedAt: sql.NullTime{
				Time:  currentTime,
				Valid: true,
			},
		}

		expected = &auth_v1.GetResponse{
			User: &auth_v1.User{
				Id: id,
				User: &auth_v1.UserInfo{
					Name:  name,
					Email: email,
					Role:  role,
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
	)

	tests := []struct {
		name              string
		args              args
		want              *auth_v1.GetResponse
		err               error
		authServerMocFunc authServerMocFunc
	}{
		{
			name: "Successful case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: expected,
			err:  nil,
			authServerMocFunc: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Get", ctx, servReq).Return(servResp, nil)

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
			authServerMocFunc: func() service.AuthService {
				mock := mocks.NewAuthService(t)
				mock.On("Get", ctx, servReq).Return(nil, serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServerMocFunc()
			api := auth.NewImplementation(authServiceMock)

			apiResp, apiErr := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, apiResp)
			require.Equal(t, tt.err, apiErr)
		})
	}

}
