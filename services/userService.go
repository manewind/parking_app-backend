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
    // Получаем информацию о пользователе, включая баланс
    query := `SELECT id, username, email, balance FROM users WHERE id = @UserID`
    err := db.QueryRow(query, sql.Named("UserID", userID)).Scan(&user.ID, &user.Username, &user.Email, &user.Balance)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return models.User{}, fmt.Errorf("пользователь с таким ID не найден")
        }
        return models.User{}, fmt.Errorf("ошибка при получении пользователя по ID: %v", err)
    }

    // Получаем все транспортные средства пользователя
    vehicleQuery := `SELECT id, license_plate, model, vehicle_type FROM vehicles WHERE user_id = @UserID`
    rows, err := db.Query(vehicleQuery, sql.Named("UserID", userID))
    
    if err != nil {
        return models.User{}, fmt.Errorf("ошибка при получении транспортных средств пользователя: %v", err)
    }
    defer rows.Close()

    // Добавляем транспортные средства в пользователя
    for rows.Next() {
        var vehicle models.Vehicle
        if err := rows.Scan(&vehicle.ID, &vehicle.LicensePlate, &vehicle.Model, &vehicle.VehicleType); err != nil {
            return models.User{}, fmt.Errorf("ошибка при сканировании транспортного средства: %v", err)
        }
        user.Vehicles = append(user.Vehicles, vehicle)
    }

    if err := rows.Err(); err != nil {
        return models.User{}, fmt.Errorf("ошибка при итерации по транспортным средствам: %v", err)
    }

    return user, nil
}


func GetAllUsers(db *sql.DB) ([]models.User, error) {
	// Словарь для сопоставления пользователей по их ID
	usersMap := make(map[int]*models.User)

	// Обновлённый SQL-запрос с добавлением поля model
	query := `SELECT 
            users.id AS user_id,
            users.username,
            users.email,
            vehicles.id AS vehicle_id,
            vehicles.license_plate,
            vehicles.vehicle_type,
            vehicles.model
        FROM 
            users
        LEFT JOIN 
            vehicles ON users.id = vehicles.user_id;`

	// Выполнение запроса
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %v", err)
	}
	defer rows.Close()

	// Перебор результатов
	for rows.Next() {
		var (
			userID       int
			username     string
			email        string
			vehicleID    sql.NullInt64
			licensePlate sql.NullString
			vehicleType  sql.NullString
			vehicleModel sql.NullString
		)

		// Чтение строки результата
		err := rows.Scan(&userID, &username, &email, &vehicleID, &licensePlate, &vehicleType, &vehicleModel)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных пользователя: %v", err)
		}

		if _, exists := usersMap[userID]; !exists {
			usersMap[userID] = &models.User{
				ID:       userID,
				Username: username,
				Email:    email,
				Vehicles: []models.Vehicle{},
			}
		}

		// Добавление информации о транспортных средствах
		if vehicleID.Valid && licensePlate.Valid && vehicleType.Valid && vehicleModel.Valid {
			vehicle := models.Vehicle{
				ID:           int(vehicleID.Int64),
				LicensePlate: licensePlate.String,
				VehicleType:  vehicleType.String,
				Model:        vehicleModel.String,
			}
			usersMap[userID].Vehicles = append(usersMap[userID].Vehicles, vehicle)
		}
	}

	// Проверка на ошибки при переборе строк
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе пользователей: %v", err)
	}

	// Преобразование карты в слайс
	users := make([]models.User, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, *user)
	}

	return users, nil
}


func DeleteUserByID(db *sql.DB, userID int) error {
    // Проверяем, существует ли пользователь с указанным ID
    var existingID int
    checkQuery := `SELECT id FROM users WHERE id = @UserID`
    err := db.QueryRow(checkQuery, sql.Named("UserID", userID)).Scan(&existingID)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("пользователь с ID %d не найден", userID)
        }
        return fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
    }

    // Удаляем пользователя
    deleteQuery := `DELETE FROM users WHERE id = @UserID`
    _, err = db.Exec(deleteQuery, sql.Named("UserID", userID))
    if err != nil {
        return fmt.Errorf("ошибка при удалении пользователя с ID %d: %v", userID, err)
    }

    return nil
}


func TopUpBalance(db *sql.DB, userID int, amount float64) error {
    // Проверяем, существует ли пользователь
    var existingBalance float64
    checkQuery := `SELECT balance FROM users WHERE id = @UserID`
    err := db.QueryRow(checkQuery, sql.Named("UserID", userID)).Scan(&existingBalance)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("пользователь с ID %d не найден", userID)
        }
        return fmt.Errorf("ошибка при проверке пользователя: %v", err)
    }

    // Пополняем баланс
    updateQuery := `UPDATE users SET balance = balance + @Amount, updated_at = CURRENT_TIMESTAMP WHERE id = @UserID`
    _, err = db.Exec(updateQuery, sql.Named("Amount", amount), sql.Named("UserID", userID))
    if err != nil {
        return fmt.Errorf("ошибка при пополнении баланса: %v", err)
    }

    return nil
}


