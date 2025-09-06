package model

import "github.com/google/uuid"

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RolePassenger UserRole = "passenger"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Firstname  string    `json:"firstname"`
	Middlename string    `json:"middlename"`
	Lastname   string    `json:"lastname"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Mobile     string    `json:"mobile"`
	Role       UserRole  `json:"role"`
}

func ParseRole(role string) UserRole {
	switch role {
	case "Admin":
		return RoleAdmin
	case "Passenger":
		return RolePassenger
	default:
		return RolePassenger
	}
}
