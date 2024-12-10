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


func GetBookingHandler(c *gin.Context) {
    // Получаем userID из параметров запроса
    userIDStr := c.Param("userID")
    fmt.Printf("Получен userID из запроса: %s\n", userIDStr)

    // Преобразуем userID в число
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        fmt.Printf("Ошибка преобразования userID: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат userID, ожидается число",
        })
        return
    }

    fmt.Printf("Преобразованный userID: %d\n", userID)

    // Подключаемся к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()
    fmt.Println("Подключение к базе данных успешно!")

    // Получаем бронирования пользователя
    bookings, err := services.GetUserBookings(dbConn, userID)
    if err != nil {
        fmt.Printf("Ошибка при получении бронирований: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при получении бронирований: %v", err),
        })
        return
    }

    // Возвращаем данные
    c.JSON(http.StatusOK, gin.H{
        "bookings": bookings,
    })
}



func GetAllBookingsHandler(c *gin.Context) {
    // Подключаемся к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()
    fmt.Println("Подключение к базе данных успешно!")

    // Инициализируем пустой массив для хранения бронирований
    bookings := []models.Booking{}

    // Формируем SQL-запрос с JOIN
    query := `
        SELECT 
            b.id, 
            b.user_id, 
            u.username, 
            b.parking_slot_id, 
            b.start_time, 
            b.end_time, 
            b.status, 
            b.created_at, 
            b.updated_at 
        FROM bookings b
        JOIN users u ON b.user_id = u.id;
    `
    fmt.Printf("Выполняется запрос: %s\n", query)

    // Выполняем запрос к базе данных
    rows, err := dbConn.Query(query)
    if err != nil {
        fmt.Printf("Ошибка выполнения запроса: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при выполнении запроса: %v", err),
        })
        return
    }
    defer rows.Close()
    fmt.Println("Запрос выполнен, начинаем чтение данных из базы")

    // Читаем данные из результата запроса
    for rows.Next() {
        var booking models.Booking
        var username string // Дополнительное поле для username
        err := rows.Scan(
            &booking.ID,
            &booking.UserID,
            &username, // Сканируем username
            &booking.ParkingSlotID,
            &booking.StartTime,
            &booking.EndTime,
            &booking.Status,
            &booking.CreatedAt,
            &booking.UpdatedAt,
        )
        if err != nil {
            fmt.Printf("Ошибка чтения данных: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Ошибка чтения данных",
            })
            return
        }
        booking.Username = username // Присваиваем username бронированию
        fmt.Printf("Считано бронирование: %+v\n", booking)
        bookings = append(bookings, booking)
    }

    // Проверяем наличие ошибок при обработке строк
    if err := rows.Err(); err != nil {
        fmt.Printf("Ошибка обработки строк результата запроса: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка обработки строк результата запроса",
        })
        return
    }

    // Логируем результат
    if len(bookings) == 0 {
        fmt.Println("Не найдено ни одного бронирования")
    } else {
        fmt.Printf("Найдено бронирований: %d\n", len(bookings))
    }

    // Возвращаем данные в ответе
    fmt.Printf("Возвращаем данные: %+v\n", bookings)
    c.JSON(http.StatusOK, gin.H{
        "bookings": bookings,
    })
}

func DeleteBookingHandler(c *gin.Context) {
    // Получаем bookingID из параметров запроса
    bookingIDStr := c.Param("bookingID")
    fmt.Printf("Получен bookingID из запроса: %s\n", bookingIDStr)

    // Преобразуем bookingID в число
    bookingID, err := strconv.Atoi(bookingIDStr)
    if err != nil {
        fmt.Printf("Ошибка преобразования bookingID: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат bookingID, ожидается число",
        })
        return
    }

    fmt.Printf("Преобразованный bookingID: %d\n", bookingID)

    // Подключаемся к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()
    fmt.Println("Подключение к базе данных успешно!")

    // Вызываем функцию удаления бронирования
    err = services.DeleteBooking(dbConn, bookingID)
    if err != nil {
        fmt.Printf("Ошибка при удалении бронирования: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при удалении бронирования: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("Бронирование с ID=%d успешно удалено", bookingID),
    })
}
