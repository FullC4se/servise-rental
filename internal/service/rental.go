package service

import (
	"Server-rental/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type RentalService interface {
	RentServer(ctx context.Context, userID, serverID int, durationHours int) (*domain.Rental, error)
	FinishRental(ctx context.Context, rentalID int) error
	GetUserRentals(ctx context.Context, userID int) ([]domain.Rental, error)
}

type rentalService struct {
	rentalRepo domain.RentalRepository
	userRepo   domain.UserRepo
	serverRepo domain.ServerRepo
}

func NewRentalService(rentalRepo domain.RentalRepository, userRepo domain.UserRepo, serverRepo domain.ServerRepo) *rentalService {
	return &rentalService{
		rentalRepo: rentalRepo,
		userRepo:   userRepo,
		serverRepo: serverRepo,
	}
}

func (s *rentalService) RentServer(ctx context.Context, userID, serverID int, durationHours int) (*domain.Rental, error) {
	if durationHours <= 0 {
		return nil, errors.New("длительность аренды должна быть больше 0 часов")
	}

	_, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с id %d не найден", userID)
		}
		return nil, fmt.Errorf("проверка пользователя: %w", err)
	}

	server, err := s.serverRepo.GetById(ctx, serverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("сервер с id %d не найден", serverID)
		}
		return nil, fmt.Errorf("проверка сервера: %w", err)
	}

	_, err = s.rentalRepo.GetActiveByServerID(ctx, serverID)
	if err == nil {
		return nil, fmt.Errorf("сервер с id %d уже арендован", serverID)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("проверка доступности сервера: %w", err)
	}

	totalAmount := server.PricePerHour * float64(durationHours)
	startedAt := time.Now()
	endedAt := startedAt.Add(time.Duration(durationHours) * time.Hour)

	rental := &domain.Rental{
		UserID:      userID,
		ServerID:    serverID,
		StartedAt:   startedAt,
		EndedAt:     endedAt,
		TotalAmount: totalAmount,
		Status:      domain.RentalActive,
	}

	err = s.rentalRepo.Create(ctx, rental)
	if err != nil {
		return nil, fmt.Errorf("создание аренды: %w", err)
	}

	err = s.serverRepo.UpdateStatus(ctx, serverID, domain.ServerRented)
	if err != nil {
		return nil, fmt.Errorf("обновление статуса сервера: %w", err)
	}

	return rental, nil
}

func (s *rentalService) FinishRental(ctx context.Context, rentalID int) error {
	rental, err := s.rentalRepo.GetById(ctx, rentalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("аренда с id %d не найдена", rentalID)
		}
		return fmt.Errorf("получение аренды: %w", err)
	}

	if rental.Status != domain.RentalActive {
		return errors.New("можно завершить только активную аренду")
	}

	err = s.rentalRepo.UpdateStatus(ctx, rentalID, domain.RentalCompleted)
	if err != nil {
		return fmt.Errorf("обновление статуса аренды: %w", err)
	}

	err = s.serverRepo.UpdateStatus(ctx, rental.ServerID, domain.ServerAvailable)
	if err != nil {
		return fmt.Errorf("обновление статуса сервера: %w", err)
	}

	return nil
}

func (s *rentalService) GetUserRentals(ctx context.Context, userID int) ([]domain.Rental, error) {
	rentals, err := s.rentalRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получение аренд пользователя: %w", err)
	}
	return rentals, nil
}
