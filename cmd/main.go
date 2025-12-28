package main

import (
	"ecom_test/internal/application"
	"ecom_test/internal/config"
)

func main() {
	cfg := config.Config{
		Addr:            ":8080",
		ShutdownTimeout: 10,
	}
	application.Run(cfg)
}
