package ginadp

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewRouter(log zerolog.Logger, gdb *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestID())
	r.Use(Logger(log))

	r.GET("/healthz", HealthHandler)
	r.GET("/health/db", HealthDBHandler(gdb))

	v1 := r.Group("/v1")
	{
		_ = v1
	}

	return r
}
