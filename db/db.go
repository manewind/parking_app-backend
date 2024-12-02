package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/denisenkom/go-mssqldb"
    "github.com/joho/godotenv"
)

// ConnectToDB открывает соединение с базой данных SQL Server.
func ConnectToDB() (*sql.DB, error) {
    // Загрузка переменных окружения из .env файла
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка при загрузке .env файла")
    }

    // Формируем строку подключения, используя переменные окружения
    connString := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&TrustServerCertificate=true",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVER"), os.Getenv("DB_NAME"))

    // Открытие соединения с базой данных
    db, err := sql.Open("sqlserver", connString)
    if err != nil {
        return nil, fmt.Errorf("ошибка подключения к базе данных: %s", err)
    }

    // Проверка подключения
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("не удалось подключиться к базе данных: %s", err)
    }

    fmt.Println("Подключение к базе данных успешно!")
    return db, nil
}
