package models

import "time"

type Booking struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    ParkingSlotID int       `json:"parking_slot_id"`
    StartTime     time.Time `json:"start_time"`
    Username       string    `json:"username"` 
    EndTime       time.Time `json:"end_time"`
    CreatedAt     time.Time `json:"created_at" `
    UpdatedAt     time.Time `json:"updated_at"`
}

type BookingRequest struct {
    UserID        int       `json:"user_id"`
    ParkingSlotID int       `json:"parking_slot_id"`
    StartTime     time.Time `json:"start_time"`
    EndTime       time.Time `json:"end_time"`
}