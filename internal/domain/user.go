package domain

import "time"

type User struct {
	Id           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	CreateUser(user User) error                 // for register
	GetUserByEmail(email string) (*User, error) // for login
	GetUserByID(id string) (*User, error)       //fetching the current user from the JWT
}
