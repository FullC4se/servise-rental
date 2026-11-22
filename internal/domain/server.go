package domain

type ServerStatus string

const (
	ServerAvailable ServerStatus = "available"
	ServerRented    ServerStatus = "rented"
)

func (server ServerStatus) IsValid() bool {
	return server == ServerAvailable || server == ServerRented
}

type Server struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	CPU          string       `json:"cpu"`
	Ram          int          `json:"ram"`
	Storage      int          `json:"storage"`
	PricePerHour float64      `json:"price_per_hour"`
	IsActive     bool         `json:"is_active"`
	Status       ServerStatus `json:"status"`
}
