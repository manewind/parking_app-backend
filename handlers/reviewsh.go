package handlers

import (
    "fmt"
    "backend/models"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "backend/services"
)

func AddReviewHandler(c *gin.Context) {
	var reviewRequest models.ReviewRequest

	if err := c.ShouldBindJSON(&reviewRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных для отзыва",
		})
		fmt.Println("Ошибка при привязке данных:", err) 
		return
	}

	// Проверка на наличие user_id
	if reviewRequest.UserID == 0 {
		fmt.Println("Ошибка: user_id не может быть равен 0")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id не может быть равен 0",
		})
		return
	}

	dbConn, err := db.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка подключения к базе данных: %v", err),
		})
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer dbConn.Close()

	// Логика создания нового отзыва
	review, err := services.CreateNewReview(dbConn, reviewRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при добавлении отзыва: %v", err),
		})
		fmt.Println("Ошибка при добавлении отзыва:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Отзыв успешно добавлен",
		"review": review,
	})
}



func GetReviewsHandler(c *gin.Context) {
    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    rows, err := dbConn.Query(`
        SELECT reviews.id, reviews.user_id, reviews.rating, reviews.comment, reviews.created_at, users.username
        FROM reviews
        JOIN users ON reviews.user_id = users.id
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при получении отзывов: %v", err),
        })
        return
    }
    defer rows.Close()

    var reviews []models.ReviewWithUser
    for rows.Next() {
        var review models.ReviewWithUser
        if err := rows.Scan(&review.ID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt, &review.Username); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": fmt.Sprintf("Ошибка при извлечении данных отзыва: %v", err),
            })
            return
        }
        reviews = append(reviews, review)
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при обработке строк результата запроса: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, reviews)
}

