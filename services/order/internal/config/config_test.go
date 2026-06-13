package config

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadRawFromEnvReadsOrderConfig(t *testing.T) {
	lookup := mapLookup(map[string]string{
		envHTTPAddr:              "127.0.0.1:8081",
		envHTTPReadHeaderTimeout: "3s",
		envHTTPShutdownTimeout:   "1s",
		envInventoryGRPCAddr:     "inventory:50052",
		envInventoryTimeout:      "4s",
		envPaymentGRPCAddr:       "payment:50051",
		envPaymentTimeout:        "5s",
		envRepoPostgresURL:       "postgres://order:secret@order-db:5432/orders?sslmode=disable",
	})

	rawConfig := loadRawFromEnv(lookup)

	require.Equal(t, "127.0.0.1:8081", rawConfig.http.addr)
	require.Equal(t, "3s", rawConfig.http.readHeaderTimeout)
	require.Equal(t, "1s", rawConfig.http.shutdownTimeout)
	require.Equal(t, "inventory:50052", rawConfig.inventory.grpcAddr)
	require.Equal(t, "4s", rawConfig.inventory.timeout)
	require.Equal(t, "payment:50051", rawConfig.payment.grpcAddr)
	require.Equal(t, "5s", rawConfig.payment.timeout)
	require.Equal(t, "postgres://order:secret@order-db:5432/orders?sslmode=disable", rawConfig.repo.postgresURL)
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv(envHTTPAddr, "127.0.0.1:8081")
	t.Setenv(envHTTPReadHeaderTimeout, "3s")
	t.Setenv(envHTTPShutdownTimeout, "1s")
	t.Setenv(envInventoryGRPCAddr, "inventory:50052")
	t.Setenv(envInventoryTimeout, "4s")
	t.Setenv(envPaymentGRPCAddr, "payment:50051")
	t.Setenv(envPaymentTimeout, "5s")
	t.Setenv(envRepoPostgresURL, "postgres://order:secret@order-db:5432/orders?sslmode=disable")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)

	require.Equal(t, "127.0.0.1:8081", cfg.HTTP().Addr())
	require.Equal(t, 3*time.Second, cfg.HTTP().ReadHeaderTimeout())
	require.Equal(t, time.Second, cfg.HTTP().ShutdownTimeout())
	require.Equal(t, "inventory:50052", cfg.Inventory().GRPCAddr())
	require.Equal(t, 4*time.Second, cfg.Inventory().Timeout())
	require.Equal(t, "payment:50051", cfg.Payment().GRPCAddr())
	require.Equal(t, 5*time.Second, cfg.Payment().Timeout())

	postgresURL, ok := cfg.Repo().PostgresURL()
	require.True(t, ok)
	require.Equal(t, "order-db:5432", postgresURL.Host)
}

func TestRawTryIntoConfigUsesDefaults(t *testing.T) {
	cfg, err := raw{}.tryIntoConfig()
	require.NoError(t, err)

	require.Equal(t, defaultHTTPAddr, cfg.HTTP().Addr())
	require.Equal(t, defaultReadHeaderTimeout, cfg.HTTP().ReadHeaderTimeout())
	require.Equal(t, defaultShutdownTimeout, cfg.HTTP().ShutdownTimeout())
	require.Equal(t, defaultInventoryAddr, cfg.Inventory().GRPCAddr())
	require.Equal(t, defaultClientTimeout, cfg.Inventory().Timeout())
	require.Equal(t, defaultPaymentAddr, cfg.Payment().GRPCAddr())
	require.Equal(t, defaultClientTimeout, cfg.Payment().Timeout())
	require.True(t, cfg.Repo().IsInMemory())

	_, ok := cfg.Repo().PostgresURL()
	require.False(t, ok)
}

func TestRawTryIntoConfigBuildsPostgresRepo(t *testing.T) {
	cfg, err := raw{
		repo: rawRepo{
			postgresURL: "postgres://order:secret@order-db:5432/orders?sslmode=disable",
		},
	}.tryIntoConfig()
	require.NoError(t, err)

	postgresURL, ok := cfg.Repo().PostgresURL()
	require.True(t, ok)
	require.False(t, cfg.Repo().IsInMemory())
	require.Equal(t, "postgres", postgresURL.Scheme)
	require.Equal(t, "order-db:5432", postgresURL.Host)
	require.Equal(t, "/orders", postgresURL.Path)
}

func TestRawTryIntoConfigAppliesOverrides(t *testing.T) {
	cfg, err := raw{
		http: rawHTTP{
			addr:              "127.0.0.1:8082",
			readHeaderTimeout: "10s",
			shutdownTimeout:   "3s",
		},
		inventory: rawInventoryClient{
			grpcAddr: "dns:///inventory:50052",
			timeout:  "1500ms",
		},
		payment: rawPaymentClient{
			grpcAddr: "payment:50051",
			timeout:  "2s",
		},
	}.tryIntoConfig()
	require.NoError(t, err)

	require.Equal(t, "127.0.0.1:8082", cfg.HTTP().Addr())
	require.Equal(t, 10*time.Second, cfg.HTTP().ReadHeaderTimeout())
	require.Equal(t, 3*time.Second, cfg.HTTP().ShutdownTimeout())
	require.Equal(t, "dns:///inventory:50052", cfg.Inventory().GRPCAddr())
	require.Equal(t, 1500*time.Millisecond, cfg.Inventory().Timeout())
	require.Equal(t, "payment:50051", cfg.Payment().GRPCAddr())
	require.Equal(t, 2*time.Second, cfg.Payment().Timeout())
}

func TestRawTryIntoConfigRejectsInvalidHTTPAddr(t *testing.T) {
	_, err := raw{
		http: rawHTTP{
			addr: "8081",
		},
	}.tryIntoConfig()

	requireFieldError(t, err, envHTTPAddr)
}

func TestRawTryIntoConfigRejectsInvalidClientAddr(t *testing.T) {
	_, err := raw{
		inventory: rawInventoryClient{
			grpcAddr: ":50052",
		},
	}.tryIntoConfig()

	requireFieldError(t, err, envInventoryGRPCAddr)
}

func TestRawTryIntoConfigRejectsURLLikeGRPCTargetWithoutEndpoint(t *testing.T) {
	_, err := raw{
		payment: rawPaymentClient{
			grpcAddr: "dns:///",
		},
	}.tryIntoConfig()

	requireFieldError(t, err, envPaymentGRPCAddr)
	require.ErrorIs(t, err, errGRPCEndpointRequired)
}

func TestRawTryIntoConfigRejectsInvalidDuration(t *testing.T) {
	_, err := raw{
		payment: rawPaymentClient{
			timeout: "-1s",
		},
	}.tryIntoConfig()

	requireFieldError(t, err, envPaymentTimeout)
	require.ErrorIs(t, err, errDurationMustBePositive)
}

func TestRawTryIntoConfigRejectsInvalidPostgresScheme(t *testing.T) {
	_, err := raw{
		repo: rawRepo{
			postgresURL: "mysql://order:secret@order-db:3306/orders",
		},
	}.tryIntoConfig()

	requireFieldError(t, err, envRepoPostgresURL)
	require.ErrorIs(t, err, errPostgresURLSchemeInvalid)
}

func TestRawTryIntoConfigRejectsPostgresURLWithoutDatabaseName(t *testing.T) {
	_, err := raw{
		repo: rawRepo{
			postgresURL: "postgres://order:secret@order-db:5432",
		},
	}.tryIntoConfig()

	requireFieldError(t, err, envRepoPostgresURL)
	require.ErrorIs(t, err, errPostgresURLDatabaseRequired)
}

func mapLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}

func requireFieldError(t *testing.T, err error, field string) {
	t.Helper()

	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), field+": "), err.Error())
	require.True(t, errors.Unwrap(err) != nil, err.Error())
}
