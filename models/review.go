package models

import "time"

type Review struct {
    ID            int       `json:"id" db:"id"`
    UserID        int       `json:"user_id" db:"user_id"`
    Rating        int       `json:"rating" db:"rating"`
    Comment       string    `json:"comment" db:"comment"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type ReviewRequest struct {
    UserID        int    `json:"user_id"`
    Rating        int    `json:"rating"`
    Comment       string `json:"comment"`
}

type ReviewWithUser struct {
    Review
    Username string `json:"username"` 
}
