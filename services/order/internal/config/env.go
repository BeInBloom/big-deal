package config

import "os"

const (
	envHTTPAddr              = "ORDER_HTTP_ADDR"
	envHTTPReadHeaderTimeout = "ORDER_HTTP_READ_HEADER_TIMEOUT"
	envHTTPShutdownTimeout   = "ORDER_HTTP_SHUTDOWN_TIMEOUT"
	envInventoryGRPCAddr     = "ORDER_INVENTORY_GRPC_ADDR"
	envInventoryTimeout      = "ORDER_INVENTORY_TIMEOUT"
	envPaymentGRPCAddr       = "ORDER_PAYMENT_GRPC_ADDR"
	envPaymentTimeout        = "ORDER_PAYMENT_TIMEOUT"
	envRepoPostgresURL       = "ORDER_REPO_POSTGRES_URL"
)

func LoadFromEnv() (Config, error) {
	raw := loadRawFromEnv(os.LookupEnv)
	return raw.tryIntoConfig()
}

func loadRawFromEnv(lookup func(string) (string, bool)) raw {
	get := func(key string) string {
		value, _ := lookup(key)
		return value
	}

	return raw{
		http: rawHTTP{
			addr:              get(envHTTPAddr),
			readHeaderTimeout: get(envHTTPReadHeaderTimeout),
			shutdownTimeout:   get(envHTTPShutdownTimeout),
		},
		inventory: rawInventoryClient{
			grpcAddr: get(envInventoryGRPCAddr),
			timeout:  get(envInventoryTimeout),
		},
		payment: rawPaymentClient{
			grpcAddr: get(envPaymentGRPCAddr),
			timeout:  get(envPaymentTimeout),
		},
		repo: rawRepo{
			postgresURL: get(envRepoPostgresURL),
		},
	}
}
