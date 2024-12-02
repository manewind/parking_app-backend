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


