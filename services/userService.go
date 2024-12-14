package services

import (
    "database/sql"
    "backend/models"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "gopkg.in/gomail.v2"
    "time"

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

func NewPayment(db *sql.DB, payment models.Payment)(models.Payment, error){

    tx, err := db.Begin()
    if err != nil {
        return models.Payment{}, fmt.Errorf("не удалось начать транзакцию: %v", err)
    }

    query := `INSERT INTO payments (user_id, amount, payment_date) OUTPUT INSERTED.id VALUES (@UserID, @Amount, CURRENT_TIMESTAMP)`
    
    var insertedID int
    err = tx.QueryRow(
        query,
        sql.Named("UserID", payment.UserID),
        sql.Named("Amount", payment.Amount),
    ).Scan(&insertedID)
    if err != nil {
        return models.Payment{}, fmt.Errorf("ошибка при вставке платежа: %v", err)
    }

    if err := tx.Commit()
    err != nil {
        return models.Payment{}, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
    }

    payment.ID = insertedID
    return payment, nil


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

    // Хешируем новый пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        fmt.Println("Ошибка при хешировании пароля:", err)
        return fmt.Errorf("ошибка при хешировании пароля: %v", err)
    }

    query := `UPDATE users SET password_hash = @newPassword, updated_at = CURRENT_TIMESTAMP WHERE email = @Email`
    fmt.Println("SQL запрос:", query)

    _, err = db.Exec(query, sql.Named("newPassword", hashedPassword), sql.Named("Email", email))
    if err != nil {
        fmt.Println("Ошибка при обновлении пароля:", err)
        return fmt.Errorf("ошибка при обновлении пароля: %v", err)
    }

    fmt.Println("Пароль успешно обновлен для email:", email)
    return nil
}


func GetUserByID(db *sql.DB, userID int) (models.User, error) {
    var user models.User

    // Получаем основную информацию о пользователе
    query := `SELECT id, username, email, balance FROM users WHERE id = @UserID`
    err := db.QueryRow(query, sql.Named("UserID", userID)).Scan(&user.ID, &user.Username, &user.Email, &user.Balance)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.User{}, fmt.Errorf("пользователь с таким ID не найден")
        }
        return models.User{}, fmt.Errorf("ошибка при получении пользователя по ID: %v", err)
    }

    // Получаем текущий активный абонемент пользователя
    membershipQuery := `
    SELECT TOP 1 id, user_id, start_date, end_date, membership_name, price, status, description, booking_hours
    FROM memberships
    WHERE user_id = @UserID AND status = 'active'
    `
    var membership models.Membership
    err = db.QueryRow(membershipQuery, sql.Named("UserID", userID)).Scan(
        &membership.ID, &membership.UserID, &membership.StartDate, &membership.EndDate,
        &membership.MembershipName, &membership.Price, &membership.Status,
        &membership.Description, &membership.BookingHours)
    if err != nil {
        if err == sql.ErrNoRows {
            user.Membership = nil // У пользователя нет активного абонемента
        } else {
            return models.User{}, fmt.Errorf("ошибка при получении абонемента: %v", err)
        }
    } else {
        user.Membership = &membership // Привязываем активный абонемент к пользователю
    }

    // Получаем транспортные средства пользователя
    vehicleQuery := `SELECT id, license_plate, model, vehicle_type FROM vehicles WHERE user_id = @UserID`
    rows, err := db.Query(vehicleQuery, sql.Named("UserID", userID))
    if err != nil {
        return models.User{}, fmt.Errorf("ошибка при получении транспортных средств пользователя: %v", err)
    }
    defer rows.Close()

    // Добавляем транспортные средства к пользователю
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
    usersMap := make(map[int]*models.User)

    query := `
        SELECT 
            users.id AS user_id, users.username, users.email,
            vehicles.id AS vehicle_id, vehicles.license_plate, vehicles.vehicle_type, vehicles.model,
            memberships.id AS membership_id, memberships.start_date, memberships.end_date, 
            memberships.membership_name, memberships.price, memberships.status,
            memberships.description, memberships.booking_hours
        FROM 
            users
        LEFT JOIN 
            vehicles ON users.id = vehicles.user_id
        LEFT JOIN 
            memberships ON users.id = memberships.user_id;
    `
    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении списка пользователей: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var (
            userID         int
            username       string
            email          string
            vehicleID      sql.NullInt64
            licensePlate   sql.NullString
            vehicleType    sql.NullString
            vehicleModel   sql.NullString
            membershipID   sql.NullInt64
            startDate      sql.NullTime
            endDate        sql.NullTime
            membershipName sql.NullString
            price          sql.NullFloat64
            status         sql.NullString
            description    sql.NullString
            bookingHours   sql.NullString
        )

        err := rows.Scan(
            &userID, &username, &email,
            &vehicleID, &licensePlate, &vehicleType, &vehicleModel,
            &membershipID, &startDate, &endDate, &membershipName,
            &price, &status, &description, &bookingHours,
        )
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

        if vehicleID.Valid {
            vehicle := models.Vehicle{
                ID:           int(vehicleID.Int64),
                LicensePlate: licensePlate.String,
                VehicleType:  vehicleType.String,
                Model:        vehicleModel.String,
            }
            usersMap[userID].Vehicles = append(usersMap[userID].Vehicles, vehicle)
        }

        if membershipID.Valid {
            membership := models.Membership{
                ID:            int(membershipID.Int64),
                UserID:        userID,
                StartDate:     startDate.Time,
                EndDate:       endDate.Time,
                MembershipName: membershipName.String,
                Price:         price.Float64,
                Status:        status.String,
                Description:   description.String,
                BookingHours:  bookingHours.String,
            }
            usersMap[userID].Membership = &membership
        }
    }

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



var jwtSecretKey = []byte("your-secret-key") 

func GenerateResetToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(), 
	})

	// Подписываем токен
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateResetToken(tokenString string) (string, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неподдерживаемый метод подписи: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("невалидный токен: %v", err)
	}

	// Проверяем, что токен действителен
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверка на истечение срока действия
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return "", fmt.Errorf("срок действия токена истек")
			}
		}
		// Возвращаем email пользователя из токена
		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
	}
	return "", fmt.Errorf("невалидный токен")
}

func SendEmail(to, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "tkacklim@gmail.com") // Здесь ваш email
    m.SetHeader("To", to) // Email получателя
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    // Используйте пароль приложения, если включена двухфакторная аутентификация
    dialer := gomail.NewDialer("smtp.gmail.com", 587, "tkacklim@gmail.com", "rxkg lhfy oupd sdew")

    // Попытка отправить письмо
    if err := dialer.DialAndSend(m); err != nil {
        return fmt.Errorf("не удалось отправить письмо: %v", err)
    }
    return nil
}




