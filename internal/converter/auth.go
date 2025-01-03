package converter

import (
	"github.com/tlb_katia/auth/internal/model"
	"github.com/tlb_katia/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) *auth_v1.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &auth_v1.User{
		Id:        user.ID,
		User:      ToUserInfoFromService(&user.User),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(user *model.UserInfo) *auth_v1.UserInfo {
	return &auth_v1.UserInfo{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            auth_v1.Role(user.Role),
	}
}

func ToUserUpdateFromService(user *model.UserUpdate) *auth_v1.UserUpdate {
	return &auth_v1.UserUpdate{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func ToUserInfoFromGRPC(user *auth_v1.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     model.Role(user.Role),
	}
}

func ToUserUpdateFromGRPC(user *auth_v1.UserUpdate) *model.UserUpdate {
	return &model.UserUpdate{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
