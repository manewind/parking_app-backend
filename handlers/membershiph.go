package handlers

import (
    "fmt"
    "backend/models"
    "backend/services"
    "backend/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func CreateMembershipHandler(c *gin.Context) {
    var membershipRequest models.Membership
    err := c.ShouldBindJSON(&membershipRequest)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат данных для создания абонемента",
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

    createdMembership, err := services.CreateMembership(dbConn, membershipRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Ошибка при создании абонемента: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, createdMembership)
}

func GetMembershipByIDHandler(c *gin.Context) {
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

    membership, err := services.GetMembershipByID(dbConn, membershipID)
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
