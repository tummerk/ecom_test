package config

import "time"

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}
