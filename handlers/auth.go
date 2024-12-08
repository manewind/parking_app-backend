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
    var loginRequest models.LoginRequest
    err := c.ShouldBindJSON(&loginRequest)
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

    storedUser, err := services.GetUserByEmail(dbConn, loginRequest.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Пользователь не найден",
        })
        return
    }

    // Логируем значения для отладки
    fmt.Println("Хэш из базы данных:", storedUser.PasswordHash)
    fmt.Println("Пароль, введенный пользователем:", loginRequest.Password)

    // Сравниваем хеш пароля
    err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(loginRequest.Password))
    if err != nil {
        fmt.Println("Ошибка сравнения паролей:", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Неверный пароль",
        })
        return
    }

    // Создаем токен
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

func ResetPasswordHandler(c *gin.Context) {
    type ResetPasswordRequest struct {
        Email       string `json:"email" binding:"required"`
        NewPassword string `json:"new_password" binding:"required"`
    }

    var req ResetPasswordRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка подключения к базе: %v", err)})
        return
    }
    defer dbConn.Close()

    _, err = services.GetUserByEmail(dbConn, req.Email)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь с таким email не найден"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
        return
    }

    if err := services.UpdatePasswordByEmail(dbConn, req.Email, string(hashedPassword)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при обновлении пароля: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно обновлён"})
}


func MeHandler(c *gin.Context) {
    userID, ok := c.Get("user_id")
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "user_id не найден в контексте",
        })
        return
    }

    userIDFloat, ok := userID.(float64)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка преобразования user_id в число",
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

    user, err := services.GetUserByID(dbConn, int(userIDFloat))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Пользователь не найден",
        })
        return
    }

    fmt.Printf("userID в контексте: %v, после преобразования: %d\n", userID, int(userIDFloat))

    isAdmin, err := services.IsAdmin(dbConn, int(userIDFloat))
    fmt.Println(isAdmin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка проверки прав администратора",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user_id":  int(userIDFloat),
        "username": user.Username,
        "email":    user.Email,
        "is_admin": isAdmin,
    })
}





