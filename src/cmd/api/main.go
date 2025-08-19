package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getenv("APP_PORT", "8080")
	r := gin.New()
	r.Use(gin.Recovery())

	// minimal health endpoint for compose
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	addr := fmt.Sprintf(":%s", port)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
