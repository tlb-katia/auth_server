package auth

import (
	"github.com/tlb_katia/auth/internal/repository"
	"github.com/tlb_katia/auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
	// TxManager
}

func NewService(repos repository.AuthRepository) service.AuthService {
	return &serv{
		authRepository: repos,
	}
}
