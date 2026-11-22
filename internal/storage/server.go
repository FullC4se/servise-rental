package storage

import (
	"Server-rental/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ServerRepo struct {
	db *sql.DB
}

func NewServerRepo(db *sql.DB) *ServerRepo {
	return &ServerRepo{db: db}
}

func (sr *ServerRepo) Create(ctx context.Context, server *domain.Server) error {

	_, err := sr.db.ExecContext(ctx, "INSERT INTO servers (name, cpu, ram, storage, price_per_hour, is_active, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		server.Name, server.CPU, server.Ram, server.Storage, server.PricePerHour, server.IsActive, server.Status)
	if err != nil {
		return fmt.Errorf("При создании сервера произошла ошибка:\n\t%v", err)
	}
	return nil
}

func (sr *ServerRepo) GetById(ctx context.Context, id int) (*domain.Server, error) {

	var server domain.Server

	err := sr.db.QueryRowContext(ctx, "SELECT id, name, cpu, ram, storage, price_per_hour, is_active, status FROM servers WHERE id = $1", id).Scan(
		&server.ID, &server.Name, &server.CPU, &server.Ram, &server.Storage, &server.PricePerHour, &server.IsActive, &server.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Сервера с id(%d) не существует, проверьте корректность данных.\nОшибка: %v", id, err)
		}
		return nil, err
	}
	return &server, nil

}

func (sr *ServerRepo) GetGoodServers(ctx context.Context) ([]domain.Server, error) {
	servers := make([]domain.Server, 0)

	rows, err := sr.db.QueryContext(ctx, "SELECT id, name, cpu, ram, storage, price_per_hour, is_active, status FROM servers WHERE status = $1 AND is_active = $2", domain.ServerAvailable, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var server domain.Server
		err = rows.Scan(&server.ID, &server.Name, &server.CPU, &server.Ram, &server.Storage, &server.PricePerHour, &server.IsActive, &server.Status)
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения строки:\n\t%v", err)
		}
		servers = append(servers, server)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Произошла ошибка во время итерации строк:\n\t%v", rows.Err())
	}
	return servers, nil
}

func (sr *ServerRepo) UpdateStatus(ctx context.Context, id int, status domain.ServerStatus) error {
	result, err := sr.db.ExecContext(ctx, "UPDATE servers SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("Ошибка обновления статуса сервера:\n\t%v", err)
	}
	codeUpdate, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if codeUpdate == 0 {
		return fmt.Errorf("сервер с id(%d) не найден", id)
	}
	return nil
}
