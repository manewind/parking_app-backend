package routes

import (
	"backend/handlers"
	"github.com/gin-gonic/gin"
)

func ParkingSlotsRoutes(r *gin.Engine) {
	r.GET("/parking-slots", handlers.GetParkingSlotsHandler)
	r.GET("/parking-slots/:slotID", handlers.GetParkingSlotByIDHandler)
	r.POST("/parking-slots", handlers.CreateParkingSlotHandler)
	r.PUT("/parking-slots/:slotID/status", handlers.UpdateParkingSlotStatusHandler)
}
