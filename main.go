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

    // Роут для получения всех пользователей
r.GET("/users", func(c *gin.Context) {
	rows, err := dbConn.Query("SELECT * FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка выполнения запроса: %v", err),
		})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка получения столбцов: %v", err),
		})
		return
	}

	for rows.Next() {
		// Создаём массив для значений
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Сканируем строку
		if err := rows.Scan(valuePtrs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Ошибка чтения строки: %v", err),
			})
			return
		}

		// Преобразуем строку в map
		rowMap := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		users = append(users, rowMap)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка обработки строк: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
})

    routes.RegisterRoutes(r)
    routes.SetupRoutes(r)
    routes.BookingRoutes(r)
	routes.ReviewRoutes(r)
	routes.BalanceRoutes(r)

    err = r.Run(":8000")
    if err != nil {
        log.Fatalf("Ошибка при запуске сервера: %v\n", err)
    }
}
