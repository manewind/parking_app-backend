package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/denisenkom/go-mssqldb"
    "github.com/joho/godotenv"
)

func ConnectToDB() (*sql.DB, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка при загрузке .env файла")
    }

    connString := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&TrustServerCertificate=true",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVER"), os.Getenv("DB_NAME"))

    db, err := sql.Open("sqlserver", connString)
    if err != nil {
        return nil, fmt.Errorf("ошибка подключения к базе данных: %s", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("не удалось подключиться к базе данных: %s", err)
    }

    fmt.Println("Подключение к базе данных успешно!")
    return db, nil
}
