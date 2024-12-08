package routes

import (
	"backend/handlers" 
	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты
func SetupRoutes(r *gin.Engine) {
	r.GET("/getAllUsers",handlers.GetAllUsersHandler)
}
