package handlers

import (
	"backend/models"
	"backend/services"
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFileHandler(c *gin.Context, db *sql.DB) {
	// Получение файла из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не передан"})
		return
	}

	// Определяем тип файла из переданного параметра
	fileType := c.PostForm("type")
	if fileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Тип файла не указан"})
		return
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось открыть файл"})
		return
	}
	defer src.Close()

	// Обрабатываем файл на основе типа
	switch fileType {
	case "payments":
		err = processPayments(src, db)
	case "bookings":
		err = processBookings(src, db)
	case "users":
		err = processUsers(src, db)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный тип файла"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки файла: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно обработан"})
}

func processPayments(fileContent io.Reader, db *sql.DB) error {
	scanner := bufio.NewScanner(fileContent)
	for scanner.Scan() {
		line := scanner.Text()
		var payment models.Payment
		_, err := fmt.Sscanf(line, "%d,%f", &payment.UserID, &payment.Amount)
		if err != nil {
			return fmt.Errorf("не удалось распарсить строку: %v", err)
		}

		_, err = services.NewPayment(db, payment)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении платежа: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	return nil
}

func processBookings(fileContent io.Reader, db *sql.DB) error {
	scanner := bufio.NewScanner(fileContent)
	for scanner.Scan() {
		line := scanner.Text()
		var booking models.Booking
		_, err := fmt.Sscanf(line, "%d,%d,%s", &booking.UserID, &booking.ParkingSlotID, &booking.StartTime)
		if err != nil {
			return fmt.Errorf("не удалось распарсить строку: %v", err)
		}

		_, err = services.CreateBooking(db, booking)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении бронирования: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	return nil
}

func processUsers(fileContent io.Reader, db *sql.DB) error {
	scanner := bufio.NewScanner(fileContent)
	for scanner.Scan() {
		line := scanner.Text()
		var user models.User
		_, err := fmt.Sscanf(line, "%d,%s,%s", &user.ID, &user.Username, &user.Email)
		if err != nil {
			return fmt.Errorf("не удалось распарсить строку: %v", err)
		}

		_, err = services.CreateUser(db, user)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении пользователя: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	return nil
}


