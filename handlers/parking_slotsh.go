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

// CreateParkingSlotHandler - хендлер для создания парковочного места
func CreateParkingSlotHandler(c *gin.Context) {
	var parkingSlot models.ParkingSlot
	err := c.ShouldBindJSON(&parkingSlot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных для создания парковочного места",
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

	// Заполнение времени создания и обновления
	parkingSlot.CreatedAt = time.Now()
	parkingSlot.UpdatedAt = time.Now()

	createdParkingSlot, err := services.CreateParkingSlot(dbConn, parkingSlot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при создании парковочного места: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, createdParkingSlot)
}

func UpdateParkingSlotStatusHandler(c *gin.Context) {
    fmt.Println("Обрабатываем запрос на обновление статуса парковочного места")
    
    slotIDStr := c.Param("slotID")
    slotID, err := strconv.Atoi(slotIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат ID парковочного места",
        })
        return
    }

    var request struct {
        IsOccupied bool `json:"is_occupied"`
    }
    err = c.ShouldBindJSON(&request)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для обновления статуса парковочного места",
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

    err = services.UpdateParkingSlotStatus(dbConn, slotID, request.IsOccupied)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при обновлении статуса парковочного места: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Статус парковочного места успешно обновлен"})
}


// GetParkingSlotsHandler - хендлер для получения всех парковочных мест
func GetParkingSlotsHandler(c *gin.Context) {
	dbConn, err := db.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
		})
		return
	}
	defer dbConn.Close()

	slots, err := services.GetParkingSlots(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при получении парковочных мест: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, slots)
}

// GetParkingSlotByIDHandler - хендлер для получения парковочного места по ID
func GetParkingSlotByIDHandler(c *gin.Context) {
	slotIDStr := c.Param("slotID")
	slotID, err := strconv.Atoi(slotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат ID парковочного места",
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

	slot, err := services.GetParkingSlotByID(dbConn, slotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при получении парковочного места: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, slot)
}
