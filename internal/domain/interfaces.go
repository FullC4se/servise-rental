package domain

import (
	"context"
)

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByNumber(ctx context.Context, number string) (*User, error)
}

type ServerRepo interface {
	Create(ctx context.Context, server *Server) error
	GetById(ctx context.Context, id int) (*Server, error)
	GetGoodServers(ctx context.Context) ([]Server, error)
	UpdateStatus(ctx context.Context, id int, status ServerStatus) error
}

type RentalRepository interface {
	Create(ctx context.Context, rental *Rental) error
	GetById(ctx context.Context, id int) (*Rental, error)
	GetByUserID(ctx context.Context, id int) ([]Rental, error)
	GetActiveByServerID(ctx context.Context, serverID int) (*Rental, error)
	UpdateStatus(ctx context.Context, id int, status RentalStatus) error
}
