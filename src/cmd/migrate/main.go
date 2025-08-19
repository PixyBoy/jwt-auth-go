package main

import (
	"log"

	"github.com/PixyBoy/jwt-auth-go/internal/adapters/db/mysql/models"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/config"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

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
		log.Fatalf("connect gorm: %v", err)
	}

	if err := gdb.AutoMigrate(
		&models.User{},
	); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	log.Println("migrations applied successfully")
}
