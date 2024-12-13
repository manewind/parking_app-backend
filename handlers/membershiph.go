package handlers

import (
    "fmt"
    "backend/models"
    "backend/services"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "encoding/json"
)

func CreateMembershipHandler(c *gin.Context) {
    var membershipRequest models.Membership

    // Считать сырые данные и сохранить для дальнейшего использования
    rawData, err := c.GetRawData()
    if err != nil {
        fmt.Printf("Ошибка чтения сырых данных: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Ошибка чтения данных запроса",
        })
        return
    }

    fmt.Printf("Полученные данные (raw): %s\n", string(rawData))

    // Привязать JSON вручную из rawData
    err = json.Unmarshal(rawData, &membershipRequest)
    if err != nil {
        fmt.Printf("Ошибка привязки JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для создания абонемента",
        })
        return
    }

    fmt.Printf("Привязанные данные: %+v\n", membershipRequest)

    // Подключение к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    // Вызов функции создания абонемента
    createdMembership, err := services.CreateMembership(dbConn, membershipRequest)
    if err != nil {
        fmt.Printf("Ошибка при создании абонемента: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при создании абонемента: %v", err),
        })
        return
    }

    fmt.Printf("Успешно созданный абонемент: %+v\n", createdMembership)
    c.JSON(http.StatusOK, createdMembership)
}



func GetMembershipByUserIDHandler(c *gin.Context) {
    // Извлекаем user_id из параметра URL
    userIDStr := c.Param("user_id")
    fmt.Printf("Полученный userID: %s\n", userIDStr)

    // Преобразуем строку в число
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат user_id, ожидается число",
        })
        return
    }

    // Подключаемся к базе данных
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    // Получаем абонемент по user_id
    membership, err := services.GetMembershipByUserID(dbConn, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": fmt.Sprintf("Ошибка при получении абонемента: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, membership)
}


func GetAllMembershipsHandler(c *gin.Context) {
    dbConn, err := db.ConnectToDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при подключении к базе данных: %v", err),
        })
        return
    }
    defer dbConn.Close()

    memberships, err := services.GetAllMemberships(dbConn)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при получении списка абонементов: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, memberships)
}

func UpdateMembershipHandler(c *gin.Context) {
    membershipIDStr := c.Param("membershipID")
    membershipID, err := strconv.Atoi(membershipIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат membershipID, ожидается число",
        })
        return
    }

    var membershipRequest models.Membership
    err = c.ShouldBindJSON(&membershipRequest)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для обновления абонемента",
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

    err = services.UpdateMembershipByID(dbConn, membershipID, membershipRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при обновлении абонемента: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Абонемент успешно обновлен",
    })
}

func DeleteMembershipHandler(c *gin.Context) {
    membershipIDStr := c.Param("membershipID")
    membershipID, err := strconv.Atoi(membershipIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат membershipID, ожидается число",
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

    err = services.DeleteMembershipByID(dbConn, membershipID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при удалении абонемента: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Абонемент успешно удален",
    })
}
