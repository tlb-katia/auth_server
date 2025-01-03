package auth

import (
	"github.com/tlb_katia/auth/internal/client/db"
	"github.com/tlb_katia/auth/internal/repository"
	"github.com/tlb_katia/auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

func NewService(repos repository.AuthRepository, tbManager db.TxManager) service.AuthService {
	return &serv{
		authRepository: repos,
		txManager:      tbManager,
	}
}
