package service

import (
	"Server-rental/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ServerService interface {
	Register(ctx context.Context, server *domain.Server, userRole domain.Role) error
	GetServerById(ctx context.Context, id int) (*domain.Server, error)
	GetGoodServers(ctx context.Context) ([]domain.Server, error)
}

type serverService struct {
	serverRepo domain.ServerRepo
}

func NewServerService(serverRepo domain.ServerRepo) *serverService {
	return &serverService{serverRepo: serverRepo}
}

func (sr *serverService) Register(ctx context.Context, server *domain.Server, userRole domain.Role) error {
	if userRole != domain.RoleAdmin {
		return errors.New("регистрировать сервера может только админ")
	}

	if server.Name == "" || server.CPU == "" || server.Ram <= 0 || server.Storage <= 0 {
		return errors.New("при попытке регистрации сервера произошла ошибка, проверьте корректность данных")
	}

	server.Status = domain.ServerAvailable
	server.IsActive = true

	err := sr.serverRepo.Create(ctx, server)
	if err != nil {
		return fmt.Errorf("при попытке создать сервер произошла ошибка: %w", err)
	}
	return nil
}

func (sr *serverService) GetServerById(ctx context.Context, id int) (*domain.Server, error) {
	server, err := sr.serverRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("сервера с id(%d) не существует: %w", id, err)
		}
		return nil, fmt.Errorf("при попытке найти сервер произошла ошибка: %w", err)
	}
	return server, nil
}

func (sr *serverService) GetGoodServers(ctx context.Context) ([]domain.Server, error) {
	servers, err := sr.serverRepo.GetGoodServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("при получении списка серверов произошла ошибка: %w", err)
	}
	if len(servers) == 0 {
		return servers, nil
	}
	return servers, nil
}
