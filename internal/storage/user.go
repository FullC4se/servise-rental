package storage

import (
	"Server-rental/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) Create(ctx context.Context, user *domain.User) error {

	_, err := ur.db.ExecContext(ctx, "INSERT INTO users (name, phone_number, email,password_hash, role) VALUES ($1, $2, $3, $4, $5)", user.Name, user.Number, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		return fmt.Errorf("При создании пользователя произошла ошибка:\n\t%w", err)
	}
	return nil
}

func (ur *UserRepo) GetById(ctx context.Context, id int) (*domain.User, error) {

	var user domain.User

	err := ur.db.QueryRowContext(ctx, "SELECT id, name, phone_number, email, password_hash, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Number, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Пользователя с id(%d) не существует, проверьте корректность данных.\nОшибка: %w", id, err)
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	err := ur.db.QueryRowContext(ctx, "SELECT id, name, phone_number, email, password_hash, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Number, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Пользователя с email(%s) не существует, проверьте корректность данных.\nОшибка: %w", email, err)
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) GetByNumber(ctx context.Context, number string) (*domain.User, error) {

	var user domain.User

	err := ur.db.QueryRowContext(ctx, "SELECT id, name, phone_number, email, password_hash, role FROM users WHERE phone_number = $1", number).Scan(&user.ID, &user.Name, &user.Number, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Пользователя с номером телефона(%s) не существует, проверьте корректность данных.\nОшибка: %w", number, err)
		}
		return nil, err
	}
	return &user, nil

}
