package app

import (
	"context"
	"github.com/tlb_katia/auth/internal/api/auth"
	"github.com/tlb_katia/auth/internal/client/db"
	"github.com/tlb_katia/auth/internal/client/db/pg"
	"github.com/tlb_katia/auth/internal/client/db/transaction"
	"github.com/tlb_katia/auth/internal/config"
	"github.com/tlb_katia/auth/internal/repository"
	authRepos "github.com/tlb_katia/auth/internal/repository/auth"
	"github.com/tlb_katia/auth/internal/service"
	authServ "github.com/tlb_katia/auth/internal/service/auth"
	"log"
)

type serviceProvider struct {
	pgConfig       config.PGConfig
	grpcConfig     config.GRPCConfig
	httpConfig     config.HTTPConfig
	dbClient       db.Client
	txManager      db.TxManager
	authRepository repository.AuthRepository
	authService    service.AuthService
	authImpl       *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}
		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatal("DBClient fail")
			return nil
		}
		return client
	}
	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	txManager := transaction.NewTransactionManager(s.DBClient(ctx).DB())
	return txManager
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		repo := authRepos.NewRepository(s.DBClient(ctx))
		return repo
	}
	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authServ.NewService(s.NoteRepository(ctx), s.TxManager(ctx))
	}
	return s.authService
}

func (s *serviceProvider) AuthImplementation(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}
	return s.authImpl
}
