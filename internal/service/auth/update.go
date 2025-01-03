package auth

import (
	"context"
	"github.com/tlb_katia/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UserUpdate) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if errTx := s.authRepository.Update(ctx, user); errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
