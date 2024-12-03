package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func ReviewRoutes(r *gin.Engine){
	r.POST("/review",handlers.AddReviewHandler)	
	r.GET("/usersReviews", handlers.GetReviewsHandler)

}