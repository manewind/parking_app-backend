package handlers

import (
    "fmt"
    "backend/models"
    "backend/services"
    "backend/db"
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

var jwtSecret = []byte("secret123")

func RegisterHandler(c *gin.Context) {
    var user models.User
    err := c.ShouldBindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных",
        })
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка при хешировании пароля",
        })
        return
    }

    user.PasswordHash = string(hash)

    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    createdUser, err := services.CreateUser(dbConn, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при создании пользователя: %v", err),
        })
        return
    }

    // Ответ с созданным пользователем
    c.JSON(http.StatusOK, createdUser)
}

func LoginHandler(c *gin.Context) {
    var user models.User
    err := c.ShouldBindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных",
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

    storedUser, err := services.GetUserByEmail(dbConn, user.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Пользователь не найден",
        })
        return
    }

    // Сравниваем хеш пароля
    err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(user.PasswordHash))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Неверный пароль",
        })
        return
    }

    // Создание JWT токена
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": storedUser.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка при создании токена",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
    })
}

func MeHandler(c *gin.Context) {
    // Получение userID из контекста
    userID, ok := c.Get("user_id")
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "user_id не найден в контексте",
        })
        return
    }

    // Преобразование userID в float64
    userIDFloat, ok := userID.(float64)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка преобразования user_id в число",
        })
        return
    }

    // Вывод значения userID
    fmt.Printf("Проверка userID: %v\n", userIDFloat)

    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    user, err := services.GetUserByID(dbConn, int(userIDFloat))
    if err != nil {
        fmt.Printf("Ошибка поиска пользователя по ID: %v\n", err)
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Пользователь не найден",
        })
        return
    }

    // Возвращение данных пользователя с логом
    fmt.Printf("Данные пользователя успешно получены для userID: %v\n", userIDFloat) // Лог успешного получения данных
    c.JSON(http.StatusOK, gin.H{
        "user_id":       int(userIDFloat), // Добавляем user_id в ответ
        "username":      user.Username,
        "email":         user.Email,
    })
}




