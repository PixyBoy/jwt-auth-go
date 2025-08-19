package ginadp

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthDBHandler(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := gdb.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
