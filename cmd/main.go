package main

import (
	"ecom_test/internal/application"
	"ecom_test/internal/config"
)

func main() { //понятно что конфиг надо задавать в env, но т.к. сторонние либы нельзя было использовать сделал просто внутри для упрощения надеюсь не сильно плохо 
	cfg := config.Config{
		Addr:            ":8080",
		ShutdownTimeout: 10,
	}
	application.Run(cfg)
}
