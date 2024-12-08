package services

import (
    "database/sql"
    "backend/models"
    "fmt"
    "time"
)

// CreateAdmin создает нового администратора
func CreateAdmin(db *sql.DB, admin models.Admin) (models.Admin, error) {
    var existingID int
    checkQuery := `SELECT id FROM admins WHERE user_id = @UserID`
    err := db.QueryRow(checkQuery, sql.Named("UserID", admin.UserID)).Scan(&existingID)
    if err != sql.ErrNoRows {
        if err == nil {
            return models.Admin{}, fmt.Errorf("администратор с таким user_id уже существует")
        }
        return models.Admin{}, fmt.Errorf("ошибка при проверке существования администратора: %v", err)
    }

    

    query := `INSERT INTO admins (user_id, email, created_at, updated_at) 
              OUTPUT INSERTED.id, INSERTED.created_at, INSERTED.updated_at 
              VALUES (@UserID, @email, @CreatedAt, @UpdatedAt)`
    var insertedID int
    var createdAt, updatedAt time.Time
    err = db.QueryRow(query,
        sql.Named("UserID", admin.UserID),
        sql.Named("email", admin.Email),
        sql.Named("CreatedAt", admin.CreatedAt),
        sql.Named("UpdatedAt", admin.UpdatedAt)).Scan(&insertedID, &createdAt, &updatedAt)

    if err != nil {
        return models.Admin{}, fmt.Errorf("ошибка при создании администратора: %v", err)
    }

    admin.ID = insertedID
    admin.CreatedAt = createdAt
    admin.UpdatedAt = updatedAt
    return admin, nil
}

func IsAdmin(db *sql.DB, userID int) (bool, error) {
    admin, err := GetAdminByUserID(db, userID)
    if err != nil {
        return false, err
    }
    return admin != nil, nil
}






// GetAdminByUserID получает администратора по UserID
func GetAdminByUserID(db *sql.DB, userID int) (*models.Admin, error) {
    var admin models.Admin
    query := `SELECT id, user_id, email, created_at, updated_at 
              FROM admins WHERE user_id = @UserID`

    err := db.QueryRow(query, sql.Named("UserID", userID)).Scan(
        &admin.ID,
        &admin.UserID,
        &admin.Email,
        &admin.CreatedAt,
        &admin.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Если администратор не найден
        }
        return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
    }

    return &admin, nil
}


// GetAdminByID получает администратора по ID
func GetAdminByID(db *sql.DB, adminID int) (models.Admin, error) {
    var admin models.Admin
    query := `SELECT id, user_id, email, created_at, updated_at 
              FROM admins WHERE id = @AdminID`

    err := db.QueryRow(query, sql.Named("AdminID", adminID)).Scan(&admin.ID, &admin.UserID, &admin.Email, &admin.CreatedAt, &admin.UpdatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return models.Admin{}, fmt.Errorf("администратор с таким ID не найден")
        }
        return models.Admin{}, fmt.Errorf("ошибка при получении администратора по ID: %v", err)
    }

    return admin, nil
}

func UpdateAdmin(db *sql.DB, admin models.Admin) (models.Admin, error) {
    query := `UPDATE admins 
              SET email = @Email, updated_at = @UpdatedAt 
              WHERE id = @ID`
    
    _, err := db.Exec(query,
        sql.Named("email", admin.Email),
        sql.Named("UpdatedAt", admin.UpdatedAt),
        sql.Named("ID", admin.ID))

    if err != nil {
        return models.Admin{}, fmt.Errorf("ошибка при обновлении администратора: %v", err)
    }

    return admin, nil
}
