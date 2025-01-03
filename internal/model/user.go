package model

import (
	"database/sql"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

type Role int32

const (
	UserRole  Role = 0
	AdminRole Role = 1
)

type UserInfo struct {
	Name            string `db:"name"`
	Email           string `db:"email"`
	Password        string `db:"-"`
	PasswordConfirm string `db:"-"`
	Role            Role   `db:"role_id"`
	RoleName        string `db:"role_name"`
}

type User struct {
	ID        int64        `db:"id"`
	User      UserInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserUpdate struct {
	ID    int64                   `db:"id"`
	Name  *wrapperspb.StringValue `db:"name"`
	Email *wrapperspb.StringValue `db:"email"`
}
