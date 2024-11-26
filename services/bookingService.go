package services

import (
    "database/sql"
    "backend/models"
    "fmt"
)

func CreateBooking(db *sql.DB, booking models.Booking) (models.Booking, error) {
    query := `INSERT INTO bookings (user_id, slot_id, start_time, end_time) OUTPUT INSERTED.id VALUES (@UserID, @SlotID, @StartTime, @EndTime)`
    var insertedID int
    err := db.QueryRow(query,
        sql.Named("UserID", booking.UserID),
        sql.Named("SlotID", booking.ParkingSlotID),
        sql.Named("StartTime", booking.StartTime),
        sql.Named("EndTime", booking.EndTime)).Scan(&insertedID)

    if err != nil {
        return models.Booking{}, fmt.Errorf("ошибка при создании бронирования: %v", err)
    }

    booking.ID = insertedID
    return booking, nil
}

