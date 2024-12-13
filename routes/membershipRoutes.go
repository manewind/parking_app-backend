package routes

import (
	"backend/handlers" 
	"github.com/gin-gonic/gin"
)

func MembershipRoutes(r *gin.Engine) {
	r.GET("/memberships", handlers.GetAllMembershipsHandler)
	r.GET("/memberships/:user_id", handlers.GetMembershipByUserIDHandler)
	r.POST("/memberships", handlers.UpdateMembershipHandler)
	r.POST("/addMembership",handlers.CreateMembershipHandler)
	r.DELETE("/memberships/:membership_id", handlers.DeleteMembershipHandler)
}
