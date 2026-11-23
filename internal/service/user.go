package service

import (
	"Server-rental/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, PasswordHash string) (*domain.User, error)
	GetUserProfile(ctx context.Context, currentUserID, targetUserID int, currentUserRole domain.Role) (*domain.User, error)
}

type userService struct {
	userRepo domain.UserRepo
}

func NewUserService(userRepo domain.UserRepo) *userService {
	return &userService{userRepo: userRepo}
}

func (us *userService) Register(ctx context.Context, user *domain.User) error {
	if user.Name == "" || user.Email == "" || user.Number == "" || user.PasswordHash == "" {
		return errors.New("при регистрации обязательные параметры не были переданы")
	}

	existingByEmail, err := us.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingByEmail != nil {
		return errors.New("пользователь с таким email уже существует")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("при попытке найти пользователя по email: %w", err)
	}

	existingByNumber, err := us.userRepo.GetByNumber(ctx, user.Number)
	if err == nil && existingByNumber != nil {
		return errors.New("пользователь с таким номером телефона уже существует")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("при попытке найти пользователя по номеру телефона: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	if err := us.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("при попытке создать пользователя: %w", err)
	}
	return nil
}

func (us *userService) Login(ctx context.Context, email, PasswordHash string) (*domain.User, error) {
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("неверный email или пароль")
		}
		return nil, fmt.Errorf("ошибка поиска пользователя: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(PasswordHash))
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	return user, nil
}

func (us *userService) GetUserProfile(ctx context.Context, currentUserID, targetUserID int, currentUserRole domain.Role) (*domain.User, error) {
	if currentUserRole != domain.RoleAdmin && currentUserID != targetUserID {
		return nil, errors.New("access denied")
	}
	user, err := us.userRepo.GetById(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}
	return user, nil
}
