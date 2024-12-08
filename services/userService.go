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

func UpdatePasswordByEmail(db *sql.DB, email, newPassword string) error {
    fmt.Println("Запрос на обновление пароля для email:", email)
    fmt.Println("Новый пароль:", newPassword)

    query := `UPDATE users SET password_hash = @newPassword, updated_at = CURRENT_TIMESTAMP WHERE email = @Email`
    
    fmt.Println("SQL запрос:", query)

    _, err := db.Exec(query, sql.Named("newPassword", newPassword), sql.Named("Email", email))
    if err != nil {
        fmt.Println("Ошибка при обновлении пароля:", err)
        return fmt.Errorf("ошибка при обновлении пароля: %v", err)
    }

    fmt.Println("Пароль успешно обновлен на:", newPassword)
    fmt.Println("Пароль успешно обновлен для email:", email)
    return nil
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

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	// Define a slice to hold the users
	var users []models.User

	// Query to select all users
	query := `SELECT id, username, email FROM users`

	// Execute the query and iterate over the rows
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %v", err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных пользователя: %v", err)
		}
		// Append the user to the users slice
		users = append(users, user)
	}

	// Check for any error that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе пользователей: %v", err)
	}

	return users, nil
}

