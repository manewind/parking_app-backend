package routes

import (
	"backend/handlers" 
	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты
func BookingRoutes(r *gin.Engine) {
	
	r.GET("/available-spots",handlers.GetAvailableSpotsHandler)
	r.POST("/booking",handlers.BookSpotHandler)
	r.GET("/user-bookings/:userID", handlers.GetUserBookingsHandler)
}
