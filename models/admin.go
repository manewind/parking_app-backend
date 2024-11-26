// models/admin.go
package models

import "time"

type Admin struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`    
    Role      string    `json:"role"`       
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
