package handlers

import (
	"backend/models"
	"backend/services"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"time"
)

func UploadFileHandler(c *gin.Context, db *sql.DB) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не передан"})
		return
	}

	fileType := c.PostForm("type")
	if fileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Тип файла не указан"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось открыть файл"})
		return
	}
	defer src.Close()

	switch fileType {
	case "payments":
		err = processPayments(src, db)
	case "bookings":
		err = processBookings(src, db)
	case "users":
		err = processUsers(src, db)
	case "reviews":
		err = processReviews(src, db)
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
	f, err := excelize.OpenReader(fileContent)
	if err != nil {
		return fmt.Errorf("не удалось открыть XLSX-файл: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("не удалось найти листы в XLSX-файле")
	}

	sheetName := sheetList[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("ошибка чтения строк с листа %s: %v", sheetName, err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			return fmt.Errorf("недостаточно данных в строке %d: %v", i+1, row)
		}

		userID, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("не удалось распарсить UserID в строке %d: %v", i+1, err)
		}
		amount, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return fmt.Errorf("не удалось распарсить Amount в строке %d: %v", i+1, err)
		}

		payment := models.Payment{UserID: userID, Amount: amount}
		_, err = services.NewPayment(db, payment)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении платежа в строке %d: %v", i+1, err)
		}
	}

	return nil
}

func processBookings(fileContent io.Reader, db *sql.DB) error {
	f, err := excelize.OpenReader(fileContent)
	if err != nil {
		return fmt.Errorf("не удалось открыть XLSX-файл: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("не удалось найти листы в XLSX-файле")
	}

	sheetName := sheetList[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("ошибка чтения строк с листа %s: %v", sheetName, err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			return fmt.Errorf("недостаточно данных в строке %d: %v", i+1, row)
		}

		userID, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("не удалось распарсить UserID в строке %d: %v", i+1, err)
		}
		parkingSlotID, err := strconv.Atoi(row[1])
		if err != nil {
			return fmt.Errorf("не удалось распарсить ParkingSlotID в строке %d: %v", i+1, err)
		}

		// Используем текущее время
		startTime := time.Now()

		booking := models.Booking{UserID: userID, ParkingSlotID: parkingSlotID, StartTime: startTime}
		_, err = services.CreateBooking(db, booking)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении бронирования в строке %d: %v", i+1, err)
		}
	}

	return nil
}


func processUsers(fileContent io.Reader, db *sql.DB) error {
	f, err := excelize.OpenReader(fileContent)
	if err != nil {
		return fmt.Errorf("не удалось открыть XLSX-файл: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("не удалось найти листы в XLSX-файле")
	}

	sheetName := sheetList[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("ошибка чтения строк с листа %s: %v", sheetName, err)
	}

	for i, row := range rows {
		if i == 0 {
			// Пропускаем заголовок
			continue
		}
		if len(row) < 3 { // Учитываем только 4 столбца: id, username, password_hash, email
			return fmt.Errorf("недостаточно данных в строке %d: %v", i+1, row)
		}

		username := row[0]
		passwordHash := row[1]
		email := row[2]

		// Создаем пользователя только с необходимыми полями
		user := models.User{
			Username:     username,
			PasswordHash: passwordHash,
			Email:        email,
		}

		_, err = services.CreateUser(db, user)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении пользователя в строке %d: %v", i+1, err)
		}
	}

	return nil
}




func processReviews(fileContent io.Reader, db *sql.DB) error {
    // Читаем файл XLSX
    f, err := excelize.OpenReader(fileContent)
    if err != nil {
        return fmt.Errorf("не удалось открыть XLSX-файл: %v", err)
    }
    defer f.Close()

    // Получаем список листов
    sheetList := f.GetSheetList()
    if len(sheetList) == 0 {
        return fmt.Errorf("не удалось найти листы в XLSX-файле")
    }

    fmt.Printf("Листы в файле: %v\n", sheetList) // Для отладки

    // Берем первый лист
    sheetName := sheetList[0]
    if sheetName == "" {
        return fmt.Errorf("не удалось определить имя первого листа")
    }

    fmt.Printf("Используемый лист: %s\n", sheetName) // Для отладки

    rows, err := f.GetRows(sheetName)
    if err != nil {
        return fmt.Errorf("ошибка чтения строк с листа %s: %v", sheetName, err)
    }

    // Обрабатываем строки
    for i, row := range rows {
        // Пропускаем заголовок (первая строка)
        if i == 0 {
            continue
        }

        if len(row) < 3 {
            return fmt.Errorf("недостаточно данных в строке %d: %v", i+1, row)
        }

        var review models.ReviewRequest

        // Парсим UserID
        userID, err := strconv.Atoi(row[0])
        if err != nil {
            return fmt.Errorf("не удалось распарсить UserID в строке %d: %v", i+1, err)
        }
        review.UserID = userID

        // Парсим Rating
        rating, err := strconv.Atoi(row[1])
        if err != nil {
            return fmt.Errorf("не удалось распарсить Rating в строке %d: %v", i+1, err)
        }
        review.Rating = rating

        // Добавляем Comment
        review.Comment = row[2]

        // Сохраняем отзыв в базу данных
        _, err = services.CreateNewReview(db, review)
        if err != nil {
            return fmt.Errorf("ошибка при сохранении отзыва в строке %d: %v", i+1, err)
        }
    }

    return nil
}





