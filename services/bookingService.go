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


func DeleteBooking(db *sql.DB, bookingID int) error {
    // Запрос на удаление бронирования по ID
    query := `DELETE FROM bookings WHERE id = @BookingID`

    // Логирование запроса
    log.Printf("Удаление бронирования с ID=%d", bookingID)

    // Выполнение запроса
    result, err := db.Exec(query, sql.Named("BookingID", bookingID))
    if err != nil {
        log.Printf("Ошибка при выполнении запроса: %s", query)
        return fmt.Errorf("ошибка при удалении бронирования: %v", err)
    }

    // Проверка, был ли удален хотя бы один ряд
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Ошибка при получении количества затронутых строк: %v", err)
        return fmt.Errorf("ошибка при получении количества затронутых строк: %v", err)
    }

    if rowsAffected == 0 {
        log.Printf("Не найдено бронирования с ID=%d для удаления", bookingID)
        return fmt.Errorf("не найдено бронирования с ID=%d для удаления", bookingID)
    }

    log.Printf("Бронирование с ID=%d успешно удалено", bookingID)
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



