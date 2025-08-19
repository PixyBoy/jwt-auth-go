package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	DB struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
		MaxOpen  int    `mapstructure:"DB_MAX_OPEN"`
		MaxIdle  int    `mapstructure:"DB_MAX_IDLE"`
	} `mapstructure:",squash"`

	Redis struct {
		Addr     string `mapstructure:"REDIS_ADDR"`
		Password string `mapstructure:"REDIS_PASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
	} `mapstructure:",squash"`

	JWT struct {
		Secret string        `mapstructure:"JWT_SECRET"`
		TTL    time.Duration `mapstructure:"-"`
		TTLMin int           `mapstructure:"JWT_TTL_MINUTES"`
	}

	OTP struct {
		Digits          int `mapstructure:"OTP_DIGITS"`
		TTLSeconds      int `mapstructure:"OTP_TTL_SECONDS"`
		MaxAttempts     int `mapstructure:"OTP_MAX_ATTEMPTS"`
		RateLimitMax    int `mapstructure:"OTP_RATE_LIMIT_MAX"`
		RateLimitWindow int `mapstructure:"OTP_RATE_LIMIT_WINDOW_SECONDS"`
	}
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("JWT_TTL_MINUTES", 1440)
	v.SetDefault("OTP_DIGITS", 6)
	v.SetDefault("OTP_TTL_SECONDS", 120)
	v.SetDefault("OTP_MAX_ATTEMPTS", 5)
	v.SetDefault("OTP_RATE_LIMIT_MAX", 3)
	v.SetDefault("OTP_RATE_LIMIT_WINDOW_SECONDS", 600)

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	c.JWT.TTL = time.Duration(c.JWT.TTLMin) * time.Minute

	return &c, nil
}
