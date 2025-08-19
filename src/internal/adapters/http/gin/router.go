package ginadp

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewRouter(log zerolog.Logger, gdb *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestID())
	r.Use(Logger(log))

	r.GET("/healthz", HealthHandler)
	r.GET("/health/db", HealthDBHandler(gdb))
	r.GET("/health/redis", HealthRedisHandler(rdb))

	v1 := r.Group("/v1")
	{
		_ = v1
	}

	return r
}
