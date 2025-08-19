package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	ginadp "github.com/PixyBoy/jwt-auth-go/internal/adapters/http/gin"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/config"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/db"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/logger"
)

type App struct {
	Cfg  *config.Config
	Log  zerolog.Logger
	DB   *gorm.DB
	HTTP *gin.Engine
}

type Engine = gin.Engine

func Build() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.AppEnv)

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

	r := ginadp.NewRouter(log, gdb)

	return &App{
		Cfg:  cfg,
		Log:  log,
		DB:   gdb,
		HTTP: r,
	}, nil
}
