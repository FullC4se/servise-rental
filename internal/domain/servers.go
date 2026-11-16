package domain

import "time"

type ServerStatus string

const (
	ServerAvailable   ServerStatus = "available"
	ServerRented      ServerStatus = "rented"
	ServerMaintenance ServerStatus = "maintenance"
)

func (server ServerStatus) IsValid() bool {
	return server == ServerAvailable || server == ServerRented || server == ServerMaintenance
}

type Server struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	CPU           string       `json:"cpu"`
	RamGB         int          `json:"ram_gb"`
	DiskGB        int          `json:"disk_gb"`
	PricePerMonth float64      `json:"price_per_month"`
	Status        ServerStatus `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
