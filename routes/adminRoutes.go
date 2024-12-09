package routes

import (
	"backend/handlers" 
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/getAllUsers",handlers.GetAllUsersHandler)
	r.GET("/user/:user_id", handlers.GetUserHandler)
	r.DELETE("/user/:user_id", handlers.DeleteUserHandler) 
	r.PUT("user/:user_id/add-balance")

}
