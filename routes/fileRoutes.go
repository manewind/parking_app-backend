package routes

import (
	"backend/handlers"
	"github.com/gin-gonic/gin"
	"database/sql"
)

func FileRoutes(r *gin.Engine, db *sql.DB) {
    r.POST("/uploadExcel", func(c *gin.Context) {
        handlers.UploadFileHandler(c, db)
    })
}
