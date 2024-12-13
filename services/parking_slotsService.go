package services

import (
    "database/sql"
    "backend/models"
    "fmt"
	"time"
    "log"
)

func CreateParkingSlot(db *sql.DB, parkingSlot models.ParkingSlot) (models.ParkingSlot, error) {
    query := `INSERT INTO parking_slots (slot_number, is_occupied, created_at, updated_at)
              OUTPUT INSERTED.id 
              VALUES (@SlotNumber, @IsOccupied, @CreatedAt, @UpdatedAt)`
    
    log.Printf("Создание парковочного места с данными: SlotNumber=%d, IsOccupied=%v, CreatedAt=%s, UpdatedAt=%s",
        parkingSlot.SlotNumber, parkingSlot.IsOccupied, parkingSlot.CreatedAt, parkingSlot.UpdatedAt)

    var insertedID int
    err := db.QueryRow(query,
        sql.Named("SlotNumber", parkingSlot.SlotNumber),
        sql.Named("IsOccupied", parkingSlot.IsOccupied),
        sql.Named("CreatedAt", parkingSlot.CreatedAt),
        sql.Named("UpdatedAt", parkingSlot.UpdatedAt)).Scan(&insertedID)

    // Логирование SQL-запроса и ошибки
    if err != nil {
        log.Printf("Ошибка при выполнении SQL-запроса: %s", query)
        log.Printf("Параметры: SlotNumber=%d, IsOccupied=%v, CreatedAt=%s, UpdatedAt=%s",
            parkingSlot.SlotNumber, parkingSlot.IsOccupied, parkingSlot.CreatedAt, parkingSlot.UpdatedAt)
        return models.ParkingSlot{}, fmt.Errorf("ошибка при создании парковочного места: %v", err)
    }

    log.Printf("Успешно создано парковочное место с ID=%d", insertedID)
    parkingSlot.ID = insertedID
    return parkingSlot, nil
}


func UpdateParkingSlotStatus(db *sql.DB, slotID int, isOccupied bool) error {
    query := `UPDATE parking_slots
              SET is_occupied = @IsOccupied, updated_at = @UpdatedAt
              WHERE id = @SlotID`

    // Логирование запроса
    log.Printf("Обновление статуса парковочного места с ID=%d на IsOccupied=%v", slotID, isOccupied)

    // Выполнение запроса
    _, err := db.Exec(query,
        sql.Named("IsOccupied", isOccupied),
        sql.Named("UpdatedAt", time.Now()),
        sql.Named("SlotID", slotID))

    if err != nil {
        log.Printf("Ошибка при выполнении запроса: %s", query)
        return fmt.Errorf("ошибка при обновлении статуса парковочного места: %v", err)
    }

    log.Printf("Статус парковочного места с ID=%d успешно обновлен", slotID)
    return nil
}

func GetParkingSlots(db *sql.DB) ([]models.ParkingSlot, error) {
    query := `
        SELECT id, slot_number, is_occupied, created_at, updated_at
        FROM parking_slots
    `

    log.Println("Получение всех парковочных мест")

    // Выполнение запроса
    rows, err := db.Query(query)
    if err != nil {
        log.Printf("Ошибка при выполнении запроса: %v", err)
        return nil, fmt.Errorf("ошибка при получении парковочных мест: %v", err)
    }
    defer rows.Close()

    var slots []models.ParkingSlot
    for rows.Next() {
        var slot models.ParkingSlot
        err := rows.Scan(&slot.ID, &slot.SlotNumber, &slot.IsOccupied, &slot.CreatedAt, &slot.UpdatedAt)
        if err != nil {
            log.Printf("Ошибка при чтении данных: %v", err)
            return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
        }
        slots = append(slots, slot)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Ошибка при обработке строк результата запроса: %v", err)
        return nil, fmt.Errorf("ошибка при обработке строк: %v", err)
    }

    return slots, nil
}


func GetParkingSlotByID(db *sql.DB, slotID int) (models.ParkingSlot, error) {
    query := `
        SELECT id, slot_number, is_occupied, created_at, updated_at
        FROM parking_slots
        WHERE id = @SlotID
    `

    log.Printf("Получение парковочного места с ID=%d", slotID)

    var slot models.ParkingSlot
    err := db.QueryRow(query, sql.Named("SlotID", slotID)).Scan(&slot.ID, &slot.SlotNumber, &slot.IsOccupied, &slot.CreatedAt, &slot.UpdatedAt)

    if err != nil {
        log.Printf("Ошибка при выполнении запроса: %v", err)
        return models.ParkingSlot{}, fmt.Errorf("ошибка при получении парковочного места: %v", err)
    }

    return slot, nil
}

