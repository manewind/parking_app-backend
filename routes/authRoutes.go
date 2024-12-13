package routes

import (
    "backend/handlers"
    "github.com/gin-gonic/gin"
    "backend/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
    r.POST("/register", handlers.RegisterHandler)
    r.POST("/login", handlers.LoginHandler)
    r.GET("/me",middleware.AuthMiddleware(),handlers.MeHandler)
    r.POST("/forgot-password",handlers.ResetPasswordHandler)
    r.POST("/reset-password",handlers.ResetPassword)
}
