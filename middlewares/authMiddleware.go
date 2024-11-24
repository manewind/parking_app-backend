package middleware

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("secret123")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		authHeader := c.GetHeader("Authorization")
		fmt.Println("asdasdasd")
		if authHeader == "" {
			fmt.Println("Ошибка: отсутствует токен авторизации") // Логирование ошибки
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует токен авторизации"})
			c.Abort()
			return
		}

		// Проверка формата токена (Bearer ...)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			fmt.Println("Ошибка: неверный формат токена") // Логирование ошибки
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		// Парсинг токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("Ошибка: неверный метод подписи") // Логирование ошибки
				return nil, fmt.Errorf("неверный метод подписи")
			}
			return jwtSecret, nil
		})

		if err != nil {
			fmt.Printf("Ошибка при парсинге токена: %v\n", err) // Логирование ошибки
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или истекший токен"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				fmt.Printf("Найден user_id: %v\n", userID) // Логирование успешного извлечения user_id
				c.Set("user_id", userID)
			} else {
				fmt.Println("Ошибка: отсутствует поле user_id в токене") // Логирование ошибки
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные в токене"})
				c.Abort()
				return
			}
		} else {
			fmt.Println("Ошибка: токен невалиден") // Логирование ошибки
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			c.Abort()
			return
		}

		c.Next()
	}
}
