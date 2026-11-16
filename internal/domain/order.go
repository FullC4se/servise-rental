package domain

import "time"

type OrderStatus string

const (
	OrderInCart    OrderStatus = "in_cart"
	OrderActive    OrderStatus = "active"
	OrderCancelled OrderStatus = "cancelled"
	OrderCompleted OrderStatus = "completed"
)

type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	ServerID  int         `json:"server_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
