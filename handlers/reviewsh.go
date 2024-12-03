package handlers

import (
    "fmt"
    "backend/models"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "database/sql"
)

func AddReviewHandler(c *gin.Context) {
    var reviewRequest models.ReviewRequest
    err := c.ShouldBindJSON(&reviewRequest)  // Привязываем данные с клиента
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для отзыва",
        })
        fmt.Println("Ошибка при привязке данных:", err)  // Логируем ошибку
        return
    }

    // Логируем полученные данные от клиента, включая user_id
    fmt.Println("Полученные данные от клиента:", reviewRequest)

    // Проверка на наличие user_id
    if reviewRequest.UserID == 0 {
        fmt.Println("Ошибка: user_id не может быть равен 0")
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "user_id не может быть равен 0",
        })
        return
    }

    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        fmt.Println("Ошибка подключения к базе данных:", err)  // Логируем ошибку подключения
        return
    }
    defer dbConn.Close()

    // Создаём новый отзыв
    review := models.Review{
        UserID:    reviewRequest.UserID,
        Rating:    reviewRequest.Rating,
        Comment:   reviewRequest.Comment,
        CreatedAt: time.Now(),
    }

    fmt.Println("Запрос на добавление отзыва:", review)  // Логируем данные отзыва

    // Вставка отзыва в базу данных
    query := `INSERT INTO reviews (user_id, rating, comment) 
				OUTPUT INSERTED.id 
				VALUES (@user_id, @rating, @comment)` // Используем именованные параметры

    fmt.Println("Запрос:", query)  // Логируем сам запрос

    var reviewID int
    err = dbConn.QueryRow(query,
        sql.Named("user_id", review.UserID),
        sql.Named("rating", review.Rating),
        sql.Named("comment", review.Comment),
    ).Scan(&reviewID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при добавлении отзыва: %v", err),
        })
        fmt.Println("Ошибка при выполнении запроса:", err)  // Логируем ошибку выполнения запроса
        return
    }

    // Возвращаем ID нового отзыва и данные
    review.ID = reviewID
    fmt.Println("Отзыв успешно добавлен, ID:", reviewID)  // Логируем успешное добавление отзыва
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

