package services

import (
    "database/sql"
    "backend/models"
    "fmt"
    "log"
)

func CreateBooking(db *sql.DB, booking models.Booking) (models.Booking, error) {
    query := `INSERT INTO bookings (user_id, parking_slot_id, start_time, end_time) 
              OUTPUT INSERTED.id 
              VALUES (@UserID, @ParkingSlotID, @StartTime, @EndTime)`
    
    // Логирование переданных значений
    log.Printf("Создание бронирования с данными: UserID=%d, ParkingSlotID=%d, StartTime=%s, EndTime=%s",
        booking.UserID, booking.ParkingSlotID, booking.StartTime, booking.EndTime)

    var insertedID int
    err := db.QueryRow(query,
        sql.Named("UserID", booking.UserID),
        sql.Named("ParkingSlotID", booking.ParkingSlotID),
        sql.Named("StartTime", booking.StartTime),
        sql.Named("EndTime", booking.EndTime)).Scan(&insertedID)

    // Логирование SQL-запроса и ошибки
    if err != nil {
        log.Printf("Ошибка при выполнении SQL-запроса: %s", query)
        log.Printf("Параметры: UserID=%d, ParkingSlotID=%d, StartTime=%s, EndTime=%s",
            booking.UserID, booking.ParkingSlotID, booking.StartTime, booking.EndTime)
        return models.Booking{}, fmt.Errorf("ошибка при создании бронирования: %v", err)
    }

    log.Printf("Успешно создано бронирование с ID=%d", insertedID)
    booking.ID = insertedID
    return booking, nil
}

func DeleteRecordByUserID(db *sql.DB, entityType string, userID int, entityID int) error {
    var query string

    // Определяем SQL-запрос в зависимости от типа сущности
    switch entityType {
    case "review":
        query = `DELETE FROM reviews WHERE user_id = @UserID AND id = @EntityID`
    case "booking":
        query = `DELETE FROM bookings WHERE user_id = @UserID AND id = @EntityID`
    default:
        return fmt.Errorf("неизвестный тип сущности: %v", entityType)
    }

    // Выполняем запрос
    _, err := db.Exec(query,
        sql.Named("UserID", userID),
        sql.Named("EntityID", entityID),
    )
    if err != nil {
        return fmt.Errorf("ошибка при удалении записи с ID %d для пользователя с ID %d: %v", entityID, userID, err)
    }

    return nil
}





func GetUserBookings(db *sql.DB, userID int) ([]models.Booking, error) {
    query := `
        SELECT 
            b.id, 
            b.user_id, 
            u.username, 
            b.parking_slot_id, 
            b.start_time, 
            b.end_time, 
            b.created_at, 
            b.updated_at 
        FROM bookings b
        JOIN users u ON b.user_id = u.id
        WHERE b.user_id = @UserID
    `

    log.Printf("Получение бронирований для пользователя с UserID=%d", userID)

    // Выполняем запрос
    rows, err := db.Query(query, sql.Named("UserID", userID))
    if err != nil {
        log.Printf("Ошибка при выполнении запроса: %v", err)
        return nil, fmt.Errorf("ошибка при получении бронирований пользователя: %v", err)
    }
    defer rows.Close()

    var bookings []models.Booking
    for rows.Next() {
        var booking models.Booking
        var username string
        err := rows.Scan(
            &booking.ID,
            &booking.UserID,
            &username, // Сканируем username
            &booking.ParkingSlotID,
            &booking.StartTime,
            &booking.EndTime,
            &booking.CreatedAt,
            &booking.UpdatedAt,
        )
        if err != nil {
            log.Printf("Ошибка при чтении данных: %v", err)
            return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
        }
        booking.Username = username // Присваиваем username бронированию
        bookings = append(bookings, booking)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Ошибка при обработке строк результата запроса: %v", err)
        return nil, fmt.Errorf("ошибка при обработке строк: %v", err)
    }

    fmt.Println(bookings)

    return bookings, nil
}





