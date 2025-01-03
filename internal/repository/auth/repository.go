package auth

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/tlb_katia/auth/internal/client/db"
	"github.com/tlb_katia/auth/internal/model"
	"github.com/tlb_katia/auth/internal/repository"
)

const (
	userTable         = "users"
	userIdColumn      = "id"
	nameColumn        = "name"
	emailColumn       = "email"
	userRoleColumn    = "role_id"
	createdTimeColumn = "created_at"
	updatedTimeColumn = "updated_at"

	roleTable    = "roles"
	roleIdColumn = "id"
	roleColumn   = "role_name"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, user *model.UserInfo) (int64, error) {
	if user.Name == "" {
		return 0, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return 0, fmt.Errorf("user email cannot be empty")
	}
	//if user.Role == 0 {
	//	return 0, fmt.Errorf("user role cannot be zero")
	//}

	query, args, err := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, userRoleColumn).
		Values(user.Name, user.Email, user.Role).
		Suffix("RETURNING id").
		ToSql()

	q := db.Query{
		Name:     "auth_repository.Create",
		QueryRaw: query,
	}

	if err != nil {
		return 0, err
	}

	var userId int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	query, args, err := sq.Select(
		fmt.Sprintf("%s.%s", userTable, userIdColumn),
		fmt.Sprintf("%s.%s", userTable, nameColumn),
		fmt.Sprintf("%s.%s", userTable, emailColumn),
		fmt.Sprintf("%s.%s", roleTable, roleColumn),
		fmt.Sprintf("%s.%s", userTable, createdTimeColumn),
		fmt.Sprintf("%s.%s", userTable, updatedTimeColumn),
	).PlaceholderFormat(sq.Dollar).
		From(userTable).
		Join(fmt.Sprintf("%s ON %s.%s = %s.%s", roleTable, userTable, userRoleColumn, roleTable, roleIdColumn)).
		Where(sq.Eq{fmt.Sprintf("%s.%s", userTable, userIdColumn): id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.Get",
		QueryRaw: query,
	}

	user := &model.User{}
	err = r.db.DB().ScanOneContext(ctx, user, q, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) Update(ctx context.Context, user *model.UserUpdate) error {
	if user.Name == nil && user.Email == nil {
		return nil
	}

	builder := sq.Update(userTable).PlaceholderFormat(sq.Dollar)

	if user.Name.Value != "" {
		builder.Set(nameColumn, user.Name.Value)
	}
	if user.Email.Value != "" {
		builder.Set(emailColumn, user.Email.Value)
	}

	query, args, err := builder.Where(sq.Eq{userIdColumn: user.ID}).ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	query, args, err := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIdColumn: id}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.Delete",
		QueryRaw: query,
	}

	fmt.Println(args)
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
