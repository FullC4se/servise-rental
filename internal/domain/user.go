package domain

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (role Role) IsValid() bool {
	return role == RoleAdmin || role == RoleUser
}

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Number       string `json:"number"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         Role   `json:"role"`
}
