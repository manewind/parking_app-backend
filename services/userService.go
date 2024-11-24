package services

import (
    "database/sql"
    "backend/models"
    "fmt"
)

func CreateUser(db *sql.DB, user models.User) (models.User, error) {
    var existingID int
    checkQuery := `SELECT id FROM users WHERE email = @Email`
    err := db.QueryRow(checkQuery, sql.Named("Email", user.Email)).Scan(&existingID)
    if err != sql.ErrNoRows {
        if err == nil {
            return models.User{}, fmt.Errorf("пользователь с таким email уже существует")
        }
        return models.User{}, fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
    }

    query := `INSERT INTO users (username, email, password_hash) OUTPUT INSERTED.id VALUES (@username, @Email, @password_hash)`
    var insertedID int
    err = db.QueryRow(query,
        sql.Named("username", user.Username),
        sql.Named("Email", user.Email),
        sql.Named("password_hash", user.PasswordHash)).Scan(&insertedID)
    
    if err != nil {
        return models.User{}, fmt.Errorf("ошибка при создании пользователя: %v", err)
    }

    user.ID = insertedID
    return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
    var user models.User
    query := `SELECT id, username, email, password_hash FROM users WHERE email = @Email`

    err := db.QueryRow(query, sql.Named("Email", email)).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return models.User{}, fmt.Errorf("пользователь с таким email не найден")
        }
        return models.User{}, fmt.Errorf("ошибка при получении пользователя: %v", err)
    }

    return user, nil
}

func GetUserByID(db *sql.DB, userID int) (models.User, error) {
    var user models.User
    query := `SELECT id, username, email FROM users WHERE id = @UserID`
    
    err := db.QueryRow(query, sql.Named("UserID", userID)).Scan(&user.ID, &user.Username, &user.Email)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return models.User{}, fmt.Errorf("пользователь с таким ID не найден")
        }
        return models.User{}, fmt.Errorf("ошибка при получении пользователя по ID: %v", err)
    }

    return user, nil
}

