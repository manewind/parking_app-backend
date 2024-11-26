package main

import (
    "fmt"
    "log"
    "backend/db"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "backend/routes"
    "net/http"
    "time"
)

func main() {
    dbConn, err := db.ConnectToDB()
    if err != nil {
        log.Fatalf("Ошибка при подключении к базе данных: %v", err)
    }
    defer dbConn.Close()

    var result int
    err = dbConn.QueryRow("SELECT 1").Scan(&result)
    if err != nil {
        log.Fatalf("Ошибка при выполнении тестового запроса: %v\n", err)
    } else {
        fmt.Printf("Соединение с базой данных установлено, результат запроса: %d\n", result)
    }

    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 
        AllowHeaders:     []string{"Content-Type", "Authorization"}, 
        AllowCredentials: true, 
        MaxAge:           24 * time.Hour, 
    }))

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    // Регистрируем все маршруты
    routes.RegisterRoutes(r)
    routes.SetupRoutes(r)
    routes.BookingRoutes(r)

    err = r.Run(":8000")
    if err != nil {
        log.Fatalf("Ошибка при запуске сервера: %v\n", err)
    }
}
