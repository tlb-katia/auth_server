package api

import (
	"github.com/tlb_katia/auth/internal/service"
	"github.com/tlb_katia/auth/pkg/auth_v1"
)

type Implementation struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(noteService service.AuthService) *Implementation {
	return &Implementation{
		authService: noteService,
	}
}
