package models

import "time"

type ParkingSlot struct {
	ID         int       `json:"id"`           
	SlotNumber int       `json:"slot_number"`  
	IsOccupied bool      `json:"is_occupied"`  
	CreatedAt  time.Time `json:"created_at"`   
	UpdatedAt  time.Time `json:"updated_at"`   
}
