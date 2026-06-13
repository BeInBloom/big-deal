package config

import "time"

const (
	defaultHTTPAddr          = ":8081"
	defaultReadHeaderTimeout = 5 * time.Second
	defaultShutdownTimeout   = 2 * time.Second
	defaultInventoryAddr     = "localhost:50052"
	defaultPaymentAddr       = "localhost:50051"
	defaultClientTimeout     = 5 * time.Second
)
