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


