package services

import (
    "database/sql"
    "backend/models"
    "fmt"
)

// Создание нового абонемента
func CreateMembership(db *sql.DB, membership models.Membership) (models.Membership, error) {
    fmt.Printf("Входящие данные для проверки существования абонемента: %+v\n", membership)

    var existingID int
    checkQuery := `SELECT id FROM memberships WHERE user_id = @UserID AND membership_name = @MembershipName`
    err := db.QueryRow(checkQuery, sql.Named("UserID", membership.UserID), sql.Named("MembershipName", membership.MembershipName)).Scan(&existingID)

    // Лог результатов проверки существования
    if err == sql.ErrNoRows {
        fmt.Println("Абонемент с таким названием для пользователя не найден.")
    } else if err != nil {
        fmt.Printf("Ошибка при проверке существования абонемента: %v\n", err)
        return models.Membership{}, fmt.Errorf("ошибка при проверке существования абонемента: %v", err)
    } else {
        fmt.Printf("Абонемент с таким названием уже существует (ID: %d)\n", existingID)
        return models.Membership{}, fmt.Errorf("пользователь уже имеет абонемент с таким названием")
    }

    // Проверка баланса пользователя
    var userBalance float64
    balanceQuery := `SELECT balance FROM users WHERE id = @UserID`
    err = db.QueryRow(balanceQuery, sql.Named("UserID", membership.UserID)).Scan(&userBalance)
    if err != nil {
        fmt.Printf("Ошибка при получении баланса пользователя: %v\n", err)
        return models.Membership{}, fmt.Errorf("ошибка при получении баланса пользователя: %v", err)
    }

    // Проверяем, достаточно ли средств для покупки абонемента
    if userBalance < membership.Price {
        fmt.Println("Недостаточно средств на балансе пользователя для создания абонемента.")
        return models.Membership{}, fmt.Errorf("недостаточно средств на балансе пользователя")
    }

    // Начало транзакции
    tx, err := db.Begin()
    if err != nil {
        fmt.Printf("Ошибка при начале транзакции: %v\n", err)
        return models.Membership{}, fmt.Errorf("ошибка при начале транзакции: %v", err)
    }

    // Обновление баланса пользователя
    updateBalanceQuery := `UPDATE users SET balance = balance - @Price WHERE id = @UserID`
    _, err = tx.Exec(updateBalanceQuery, sql.Named("Price", membership.Price), sql.Named("UserID", membership.UserID))
    if err != nil {
        tx.Rollback()
        fmt.Printf("Ошибка при обновлении баланса пользователя: %v\n", err)
        return models.Membership{}, fmt.Errorf("ошибка при обновлении баланса пользователя: %v", err)
    }

    // Вставка нового абонемента
    query := `INSERT INTO memberships (user_id, start_date, end_date, membership_name, price, status, description, booking_hours) 
              OUTPUT INSERTED.id VALUES (@UserID, @StartDate, @EndDate, @MembershipName, @Price, @Status, @Description, @BookingHours)`
    var insertedID int
    err = tx.QueryRow(query,
        sql.Named("UserID", membership.UserID),
        sql.Named("StartDate", membership.StartDate),
        sql.Named("EndDate", membership.EndDate),
        sql.Named("MembershipName", membership.MembershipName),
        sql.Named("Price", membership.Price),
        sql.Named("Status", membership.Status),
        sql.Named("Description", membership.Description),
        sql.Named("BookingHours", membership.BookingHours)).Scan(&insertedID)

    if err != nil {
        tx.Rollback()
        fmt.Printf("Ошибка при создании абонемента: %v\n", err)
        return models.Membership{}, fmt.Errorf("ошибка при создании абонемента: %v", err)
    }

    // Завершаем транзакцию
    err = tx.Commit()
    if err != nil {
        fmt.Printf("Ошибка при коммите транзакции: %v\n", err)
        return models.Membership{}, fmt.Errorf("ошибка при коммите транзакции: %v", err)
    }

    membership.ID = insertedID
    fmt.Printf("Успешно вставленный абонемент: %+v\n", membership)
    return membership, nil
}



// Получить абонемент по ID
func GetMembershipByUserID(db *sql.DB, userID int) (models.Membership, error) {
    var membership models.Membership
    // Обновляем запрос, чтобы искать по user_id
    query := `SELECT id, user_id, start_date, end_date, membership_name, price, status, description, booking_hours
              FROM memberships WHERE user_id = @UserID`

    // Выполняем запрос и заполняем структуру
    err := db.QueryRow(query, sql.Named("UserID", userID)).Scan(
        &membership.ID, &membership.UserID, &membership.StartDate, &membership.EndDate,
        &membership.MembershipName, &membership.Price, &membership.Status,
        &membership.Description, &membership.BookingHours)

    if err != nil {
        if err == sql.ErrNoRows {
            return models.Membership{}, fmt.Errorf("абонемент для пользователя с таким ID не найден")
        }
        return models.Membership{}, fmt.Errorf("ошибка при получении абонемента по user_id: %v", err)
    }

    return membership, nil
}


// Получить все абонементы
func GetAllMemberships(db *sql.DB) ([]models.Membership, error) {
    var memberships []models.Membership
    query := `SELECT id, user_id, start_date, end_date, membership_name, price, status, description, booking_hours FROM memberships`

    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении списка абонементов: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var membership models.Membership
        if err := rows.Scan(&membership.ID, &membership.UserID, &membership.StartDate, &membership.EndDate,
            &membership.MembershipName, &membership.Price, &membership.Status,
            &membership.Description, &membership.BookingHours); err != nil {
            return nil, fmt.Errorf("ошибка при сканировании абонемента: %v", err)
        }
        memberships = append(memberships, membership)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка при переборе абонементов: %v", err)
    }

    return memberships, nil
}

// Обновить абонемент по ID
func UpdateMembershipByID(db *sql.DB, membershipID int, newMembership models.Membership) error {
    query := `UPDATE memberships SET start_date = @StartDate, end_date = @EndDate, 
              membership_name = @MembershipName, price = @Price, status = @Status, 
              description = @Description, booking_hours = @BookingHours 
              WHERE id = @MembershipID`

    _, err := db.Exec(query,
        sql.Named("StartDate", newMembership.StartDate),
        sql.Named("EndDate", newMembership.EndDate),
        sql.Named("MembershipName", newMembership.MembershipName),
        sql.Named("Price", newMembership.Price),
        sql.Named("Status", newMembership.Status),
        sql.Named("Description", newMembership.Description),
        sql.Named("BookingHours", newMembership.BookingHours),
        sql.Named("MembershipID", membershipID))

    if err != nil {
        return fmt.Errorf("ошибка при обновлении абонемента: %v", err)
    }

    return nil
}

// Удалить абонемент по ID
func DeleteMembershipByID(db *sql.DB, membershipID int) error {
    // Проверяем, существует ли абонемент с указанным ID
    var existingID int
    checkQuery := `SELECT id FROM memberships WHERE id = @MembershipID`
    err := db.QueryRow(checkQuery, sql.Named("MembershipID", membershipID)).Scan(&existingID)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("абонемент с ID %d не найден", membershipID)
        }
        return fmt.Errorf("ошибка при проверке существования абонемента: %v", err)
    }

    // Удаляем абонемент
    deleteQuery := `DELETE FROM memberships WHERE id = @MembershipID`
    _, err = db.Exec(deleteQuery, sql.Named("MembershipID", membershipID))
    if err != nil {
        return fmt.Errorf("ошибка при удалении абонемента с ID %d: %v", membershipID, err)
    }

    return nil
}
