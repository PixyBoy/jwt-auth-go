package ginadp

import (
	_ "github.com/PixyBoy/jwt-auth-go/docs"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRouter(log zerolog.Logger, gdb *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestID())
	r.Use(Logger(log))

	// health
	r.GET("/healthz", HealthHandler)
	r.GET("/health/db", HealthDBHandler(gdb))
	r.GET("/health/redis", HealthRedisHandler(rdb))

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// v1
	v1 := r.Group("/v1")
	{
		_ = v1
	}

	return r
}

func GroupV1(r *gin.Engine) *gin.RouterGroup {
	return r.Group("/v1")
}
