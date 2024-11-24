// models/admin.go
package models

import "time"

type Admin struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`    // Ссылка на пользователя
    Role      string    `json:"role"`       // Роль администратора (например, "admin")
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
