package repository

import "github.com/tlb_katia/auth/internal/client/db"

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) {

}
