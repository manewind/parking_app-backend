package handlers

import (
	"fmt"
	"net/http"
	"backend/models"
	"backend/services"
	"backend/db" // импорт вашей базы данных
	"time"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreateAdminHandler(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Ошибка при обработке запроса: %v", err),
		})
		return
	}

	// Устанавливаем дату создания и обновления
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	// Получаем подключение к базе данных
	dbConn, err := db.ConnectToDB() // Используем GetDB вместо ConnectToDB
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
		})
		return
	}

	// Создаем администратора через сервис
	createdAdmin, err := services.CreateAdmin(dbConn, admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при создании администратора: %v", err),
		})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, createdAdmin)
}

func GetAdminByUserIDHandler(c *gin.Context) {
	userID := c.Param("user_id")

	// Преобразуем userID из string в int
	userIDInt, err := strconv.Atoi(userID) // Используем strconv.Atoi для конвертации строки в целое число
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Ошибка при преобразовании user_id в число: %v", err),
		})
		return
	}

	// Получаем подключение к базе данных
	dbConn, err := db.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
		})
		return
	}

	// Получаем администратора по user_id
	admin, err := services.GetAdminByUserID(dbConn, userIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Ошибка при получении администратора: %v", err),
		})
		return
	}

	// Отправляем данные администратора
	c.JSON(http.StatusOK, admin)
}

func UpdateAdminHandler(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Ошибка при обработке запроса: %v", err),
		})
		return
	}

	// Устанавливаем дату обновления
	admin.UpdatedAt = time.Now()

	// Получаем подключение к базе данных
	dbConn, err := db.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
		})
		return
	}

	// Обновляем данные администратора через сервис
	updatedAdmin, err := services.UpdateAdmin(dbConn, admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при обновлении администратора: %v", err),
		})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, updatedAdmin)
}
