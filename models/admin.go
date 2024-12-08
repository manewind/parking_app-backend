// models/admin.go
package models

import "time"

type Admin struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"` 
    Email       string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
