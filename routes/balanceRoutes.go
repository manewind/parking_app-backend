package routes

import (
	"backend/handlers" 
	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты
func BalanceRoutes(r *gin.Engine) {
	r.POST("/add-balance",handlers.TopUpBalanceHandler)
	r.POST("/newPayment",handlers.NewPayment)
}
