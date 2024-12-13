package models

import "time"

type User struct {
	ID           int           `json:"id"`
	Username     string        `json:"username"`
	PasswordHash string        `json:"password_hash"`
	Email        string        `json:"email"`
	Balance      float64       `json:"balance"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Vehicles     []Vehicle     `json:"vehicles"`
	Membership   *Membership   `json:"membership"` 
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
