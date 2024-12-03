package handlers

import (
    "fmt"
    "backend/models"
    "backend/services"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "strconv"
)

func BookSpotHandler(c *gin.Context) {
    var bookingRequest models.BookingRequest
    err := c.ShouldBindJSON(&bookingRequest)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для бронирования",
        })
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    booking := models.Booking{
        UserID:        bookingRequest.UserID,
        ParkingSlotID: bookingRequest.ParkingSlotID,
        StartTime:     bookingRequest.StartTime,
        EndTime:       bookingRequest.EndTime,
        Status:        "pending",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    createdBooking, err := services.CreateBooking(dbConn, booking)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при создании бронирования: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, createdBooking)
}

func GetAvailableSpotsHandler(c *gin.Context) {
    startTime := c.Query("start_time")
    endTime := c.Query("end_time")

    if startTime == "" || endTime == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing start_time or end_time"})
        return
    }

    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err)})
        return
    }
    defer dbConn.Close()

    query := `
    SELECT DISTINCT spot_id 
    FROM parking_spots 
    WHERE spot_id NOT IN (
        SELECT spot_id 
        FROM bookings 
        WHERE (start_time, end_time) OVERLAPS ($1, $2)
    )`
    rows, err := dbConn.Query(query, startTime, endTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    defer rows.Close()

    var availableSpots []int
    for rows.Next() {
        var spotID int
        if err := rows.Scan(&spotID); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading data"})
            return
        }
        availableSpots = append(availableSpots, spotID)
    }

    c.JSON(http.StatusOK, gin.H{
        "available_spots": availableSpots,
    })
}


func GetUserBookingsHandler(c *gin.Context) {
    userIDStr := c.Param("userID")
    fmt.Printf("Получен userID из запроса: %s\n", userIDStr)

    // Преобразуем userID из строки в целое число
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат userID, ожидается число",
        })
        return
    }

    fmt.Printf("Преобразованный userID: %d\n", userID)

    // Подключаемся к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    var bookingRequests []models.BookingRequest

    // Запрос к базе данных для получения нужных данных
    query := "SELECT user_id, parking_slot_id, start_time, end_time FROM bookings WHERE user_id = $1"
    fmt.Printf("Выполняется запрос: %s с параметром %d\n", query, userID)
    rows, err := dbConn.Query(query, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при получении бронирований: %v", err),
        })
        return
    }
    defer rows.Close()

    fmt.Println("Запрос выполнен, начинаем чтение данных из базы")
    // Чтение данных из результата запроса
    for rows.Next() {
        var bookingRequest models.BookingRequest
        err := rows.Scan(&bookingRequest.UserID, &bookingRequest.ParkingSlotID, &bookingRequest.StartTime, &bookingRequest.EndTime)
        if err != nil {
            fmt.Printf("Ошибка чтения данных: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
            return
        }
        fmt.Printf("Считано бронирование: %+v\n", bookingRequest)
        bookingRequests = append(bookingRequests, bookingRequest)
    }

    if len(bookingRequests) == 0 {
        fmt.Println("Не найдено бронирований для данного пользователя")
    }

    // Возвращаем данные в ответе
    fmt.Printf("Возвращаем данные: %+v\n", bookingRequests)
    c.JSON(http.StatusOK, gin.H{
        "userID":   userID,
        "bookings": bookingRequests,
    })
}
