package auth

import (
	"context"
	"github.com/tlb_katia/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, user *model.UserInfo) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		id2, errTx := s.authRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.authRepository.Get(ctx, id2)
		if errTx != nil {
			return errTx
		}
		id = id2
		return nil
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}
