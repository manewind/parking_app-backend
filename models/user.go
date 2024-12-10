package models

import "time"

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	PasswordHash string   `json:"password_hash"`
	Email       string    `json:"email"`
	Balance     float64		  `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Vehicles    []Vehicle `json:"vehicles"` 
}


type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}
