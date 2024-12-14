package handlers

import (
    "fmt"
    "backend/services"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "backend/models"
)


func TopUpBalanceHandler(c *gin.Context) {
    var request struct {
        UserID int     `json:"user_id"`
        Amount float64 `json:"amount"`
    }

    // Парсинг JSON-запроса
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для пополнения баланса",
        })
        return
    }

    // Проверка данных
    if request.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Сумма пополнения должна быть больше 0",
        })
        return
    }

    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    // Пополнение баланса
    err = services.TopUpBalance(dbConn, request.UserID, request.Amount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при пополнении баланса: %v", err),
        })
        return
    }


    c.JSON(http.StatusOK, gin.H{
        "message": "Баланс успешно пополнен",
    })
}


func NewPayment(c *gin.Context) {
    var request struct {
        UserID int     `json:"user_id"`
        Amount float64 `json:"amount"`
    }

    // Парсинг JSON-запроса
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для добавления платежа",
        })
        return
    }

    // Проверка данных
    if request.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Сумма платежа должна быть больше 0",
        })
        return
    }

    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    // Создание нового платежа
    payment := models.Payment{
        UserID: request.UserID,
        Amount: request.Amount,
    }

    newPayment, err := services.NewPayment(dbConn, payment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при создании платежа: %v", err),
        })
        return
    }

    // Возвращаем успешный ответ
    c.JSON(http.StatusOK, gin.H{
        "message": "Платеж успешно добавлен",
        "payment": newPayment,
    })
}

