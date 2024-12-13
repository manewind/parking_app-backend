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
    var registerRequest models.RegisterRequest
    err := c.ShouldBindJSON(&registerRequest)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных",
        })
        return
    }

    // Логируем пароль до хеширования
    fmt.Println("Пароль перед хешированием:", registerRequest.Password)

    // Хешируем пароль
    hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка при хешировании пароля",
        })
        return
    }

    // Создаем пользователя, сохраняя хэшированный пароль
    user := models.User{
        Username:     registerRequest.Username,
        Email:        registerRequest.Email,
        PasswordHash: string(hash), // Сохраняем хэшированный пароль
        Balance:      0,            // Начальный баланс
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

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

    c.JSON(http.StatusOK, createdUser)
}



func LoginHandler(c *gin.Context) {
    var loginRequest models.LoginRequest

    // Логируем входящий запрос
    fmt.Println("Получен запрос на логин:", c.Request.Method, c.Request.URL.Path)

    // Привязка JSON данных
    err := c.ShouldBindJSON(&loginRequest)
    if err != nil {
        fmt.Println("Ошибка привязки JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных",
        })
        return
    }

    // Логируем данные из запроса после привязки
    fmt.Println("Данные из запроса:", loginRequest)

    dbConn, err := db.ConnectToDB()
    if err != nil {
        fmt.Println("Ошибка подключения к базе данных:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer func() {
        fmt.Println("Закрытие соединения с базой данных.")
        dbConn.Close()
    }()

    // Логируем перед запросом пользователя из базы
    fmt.Println("Поиск пользователя по email:", loginRequest.Email)
    storedUser, err := services.GetUserByEmail(dbConn, loginRequest.Email)
    if err != nil {
        fmt.Println("Ошибка при получении пользователя из базы:", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Пользователь не найден",
        })
        return
    }
    fmt.Println("Найден пользователь в базе данных:", storedUser)

    // Логируем перед сравнением паролей
    fmt.Println("Пароль, введённый пользователем:", loginRequest.Password)
    fmt.Println("Хэш пароля из базы данных:", storedUser.PasswordHash)

    // Сравнение пароля с хешем
    err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(loginRequest.Password))
    if err != nil {
        fmt.Println("Ошибка сравнения паролей:", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Неверный пароль",
        })
        return
    }

    fmt.Println("Пароль успешно проверен.")

    // Создаем JWT токен
    fmt.Println("Создание JWT токена для пользователя ID:", storedUser.ID)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": storedUser.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        fmt.Println("Ошибка при создании токена:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка при создании токена",
        })
        return
    }
    fmt.Println("JWT токен успешно создан:", tokenString)

    // Отправляем ответ
    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
    })
}




func ResetPasswordHandler(c *gin.Context) {
    type ResetPasswordRequest struct {
        Email string `json:"email" binding:"required"`
    }

    var req ResetPasswordRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
        return
    }

    // Проверяем, существует ли пользователь с данным email
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
        return
    }
    defer dbConn.Close()

    _, err = services.GetUserByEmail(dbConn, req.Email)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь с таким email не найден"})
        return
    }

    resetToken, err := services.GenerateResetToken(req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
        return
    }
    resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
    err = services.SendEmail(req.Email, "Сброс пароля", fmt.Sprintf("Перейдите по ссылке для сброса пароля: %s", resetLink))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при отправке письма: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Инструкция по сбросу пароля отправлена на указанный email"})
}

func ResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Проверяем токен
	email, err := services.ValidateResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Обновляем пароль пользователя в базе данных
	dbConn, err := db.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer dbConn.Close()

	// Обновляем пароль в базе данных
	err = services.UpdatePasswordByEmail(dbConn, email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пароля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно сброшен"})
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





