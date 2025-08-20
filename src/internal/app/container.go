package app

import (
	ginadp "github.com/PixyBoy/jwt-auth-go/internal/adapters/http/gin"
	jwtadp "github.com/PixyBoy/jwt-auth-go/internal/adapters/token/jwt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	redisadp "github.com/PixyBoy/jwt-auth-go/internal/adapters/cache/redis"
	mysqladp "github.com/PixyBoy/jwt-auth-go/internal/adapters/db/mysql"
	"github.com/PixyBoy/jwt-auth-go/internal/core/services"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/config"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/db"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/logger"
)

type App struct {
	Cfg         *config.Config
	Log         zerolog.Logger
	DB          *gorm.DB
	RDB         *redis.Client
	AuthService services.AuthService
	HTTP        *Engine
}

type Engine = gin.Engine

func Build() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.AppEnv)

	// DB
	mysqlCfg := db.MySQLConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		MaxOpen:  cfg.DB.MaxOpen,
		MaxIdle:  cfg.DB.MaxIdle,
	}
	gdb, err := db.NewGorm(mysqlCfg, cfg.AppEnv)
	if err != nil {
		return nil, err
	}

	// Redis
	rdb := redisadp.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// Adapters
	otpStore := redisadp.NewOTPStore(rdb, "")
	rateLimiter := redisadp.NewRateLimiter(rdb, "")
	userRepo := mysqladp.NewUserRepo(gdb)

	// issuer
	jwtIssuer := jwtadp.NewIssuerHS256(cfg.JWT.Secret)

	// Services
	authSvc := services.NewAuthService(
		otpStore, rateLimiter, userRepo, log,
		cfg.OTP.Digits,
		cfg.OTP.TTLSeconds,
		cfg.OTP.MaxAttempts,
		cfg.OTP.RateLimitMax,
		cfg.OTP.RateLimitWindow,
		cfg.JWT.TTL,
	)
	// Router
	r := ginadp.NewRouter(log, gdb, rdb)

	v1 := ginadp.GroupV1(r)
	v1.POST("/auth/otp/request", ginadp.RequestOTPHandler(authSvc, log))
	v1.POST("/auth/otp/verify", ginadp.VerifyOTPHandler(authSvc, log))

	// protected group
	authz := ginadp.Authz(jwtIssuer)
	pg := v1.Group("/")
	pg.Use(authz)
	pg.GET("users/me", ginadp.GetMeHandler(userRepo))
	pg.GET("users/:id", ginadp.GetUserByIDHandler(userRepo))

	return &App{
		Cfg:         cfg,
		Log:         log,
		DB:          gdb,
		RDB:         rdb,
		AuthService: authSvc,
		HTTP:        r,
	}, nil
}
