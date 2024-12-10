package models

import "time"

type Vehicle struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"` 
	LicensePlate string    `json:"license_plate"`
	VehicleType  string    `json:"vehicle_type"`
	Model        string    `json:"model"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}