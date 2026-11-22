package domain

import "time"

type RentalStatus string

const (
	RentalActive    RentalStatus = "active"
	RentalCancelled RentalStatus = "cancelled"
	RentalCompleted RentalStatus = "completed"
)

type Rental struct {
	ID          int          `json:"id"`
	UserID      int          `json:"user_id"`
	ServerID    int          `json:"server_id"`
	StartedAt   time.Time    `json:"started_at"`
	EndedAt     time.Time    `json:"ended_at"`
	TotalAmount float64      `json:"total_amount"`
	Status      RentalStatus `json:"status"`
}
