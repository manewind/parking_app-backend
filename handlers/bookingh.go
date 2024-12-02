package handlers

import (
    "fmt"
    "backend/models"
    "backend/services"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
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
	userID := c.Param("userID")

	dbConn, err := db.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
		})
		return
	}
	defer dbConn.Close()

	var bookings []models.Booking
	query := "SELECT * FROM bookings WHERE user_id = $1"
	rows, err := dbConn.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при получении бронирований: %v", err),
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.ParkingSlotID, &booking.StartTime, &booking.EndTime, &booking.Status, &booking.CreatedAt, &booking.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, gin.H{
        "userID":   userID,
        "bookings": bookings,
    })
}

