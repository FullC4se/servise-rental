package storage

import (
	"Server-rental/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RentalRepo struct {
	db *sql.DB
}

func NewRentalRepo(db *sql.DB) *RentalRepo {
	return &RentalRepo{db: db}
}

func (r *RentalRepo) Create(ctx context.Context, rental *domain.Rental) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO rentals (user_id, server_id, started_at, ended_at, total_amount, status) VALUES ($1, $2, $3, $4, $5, $6)",
		rental.UserID, rental.ServerID, rental.StartedAt, rental.EndedAt, rental.TotalAmount, rental.Status,
	)
	if err != nil {
		return fmt.Errorf("При организации аренды произошла ошибка:\n\t%v", err)
	}
	return nil
}

func (r *RentalRepo) GetById(ctx context.Context, id int) (*domain.Rental, error) {

	var rental domain.Rental

	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, server_id, started_at,ended_at, total_amount,status FROM rentals WHERE id = $1", id).Scan(
		&rental.ID, &rental.UserID, &rental.ServerID, &rental.StartedAt, &rental.EndedAt, &rental.TotalAmount, &rental.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Аренды с id(%d) не существует, проверьте корректность данных.\nОшибка: %v", id, err)
		}
		return nil, err
	}
	return &rental, nil
}

func (r *RentalRepo) GetByUserID(ctx context.Context, userID int) ([]domain.Rental, error) {
	rentals := make([]domain.Rental, 0)

	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, server_id, started_at, ended_at, total_amount, status FROM rentals WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rental domain.Rental
		err = rows.Scan(&rental.ID, &rental.UserID, &rental.ServerID, &rental.StartedAt, &rental.EndedAt, &rental.TotalAmount, &rental.Status)
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения строки:\n\t%v", err)
		}
		rentals = append(rentals, rental)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("Произошла ошибка во время итерации строк:\n\t%v", rows.Err())
	}
	return rentals, nil
}

func (r *RentalRepo) GetActiveByServerID(ctx context.Context, serverID int) (*domain.Rental, error) {
	var rental domain.Rental

	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, server_id, started_at, ended_at, total_amount, status FROM rentals WHERE server_id = $1 AND status = $2",
		serverID, domain.RentalActive,
	).Scan(&rental.ID, &rental.UserID, &rental.ServerID, &rental.StartedAt, &rental.EndedAt, &rental.TotalAmount, &rental.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("Ошибка получения строки аренды: %v", err)
	}
	return &rental, nil
}

func (r *RentalRepo) UpdateStatus(ctx context.Context, id int, status domain.RentalStatus) error {
	result, err := r.db.ExecContext(ctx, "UPDATE rentals SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("При попытке изменить статус сервера произошла ошибка: %v", err)
	}

	codeUpdate, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Ошибка получения статуса запроса: %v", err)
	}
	if codeUpdate == 0 {
		return fmt.Errorf("Аренда с id(%d) не найден", id)
	}
	return nil
}
