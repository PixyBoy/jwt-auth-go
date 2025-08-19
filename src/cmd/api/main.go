package main

import (
	"fmt"
	"log"

	"github.com/PixyBoy/jwt-auth-go/internal/app"
)

func main() {
	a, err := app.Build()
	if err != nil {
		log.Fatalf("boot error: %v", err)
	}

	addr := fmt.Sprintf(":%s", a.Cfg.AppPort)
	if err := a.HTTP.Run(addr); err != nil {
		log.Fatalf("http error: %v", err)
	}
}
