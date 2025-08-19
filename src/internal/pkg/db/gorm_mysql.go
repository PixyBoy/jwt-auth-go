package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	MaxOpen  int
	MaxIdle  int
}

func dsn(c MySQLConfig) string {
	// parseTime=true
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&charset=utf8mb4,utf8",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func NewGorm(c MySQLConfig, appEnv string) (*gorm.DB, error) {
	var level glogger.LogLevel = glogger.Silent
	if appEnv == "development" {
		level = glogger.Info
	}
	gormLogger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	fmt.Printf("Connecting to DB %s:%d user=%s db=%s\n", c.Host, c.Port, c.User, c.Name)
	db, err := gorm.Open(mysql.Open(dsn(c)), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if c.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpen)
	}
	if c.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdle)
	}
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
