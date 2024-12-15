package routes

import (
	"github.com/gin-gonic/gin"
)

func FileRoutes(r *gin.Engine) {
	r.POST("/uploadExcel")
}
