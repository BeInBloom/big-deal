package config

import (
	"net/url"
	"time"
)

type (
	Config interface {
		HTTP() HTTP
		Inventory() InventoryClient
		Payment() PaymentClient
		Repo() Repo

		isConfig()
	}

	HTTP interface {
		Addr() string
		ReadHeaderTimeout() time.Duration
		ShutdownTimeout() time.Duration

		isHTTP()
	}

	InventoryClient interface {
		GRPCAddr() string
		Timeout() time.Duration

		isInventoryClient()
	}

	PaymentClient interface {
		GRPCAddr() string
		Timeout() time.Duration

		isPaymentClient()
	}

	Repo interface {
		IsInMemory() bool
		PostgresURL() (url.URL, bool)

		isRepo()
	}

	config struct {
		http      httpConfig
		inventory inventoryClient
		payment   paymentClient
		repo      repoConfig
	}

	httpConfig struct {
		addr              string
		readHeaderTimeout time.Duration
		shutdownTimeout   time.Duration
	}

	inventoryClient struct {
		grpcAddr string
		timeout  time.Duration
	}

	paymentClient struct {
		grpcAddr string
		timeout  time.Duration
	}

	repoConfig struct {
		postgresURL *url.URL
	}
)

type (
	raw struct {
		http      rawHTTP
		inventory rawInventoryClient
		payment   rawPaymentClient
		repo      rawRepo
	}

	rawHTTP struct {
		addr              string
		readHeaderTimeout string
		shutdownTimeout   string
	}

	rawInventoryClient struct {
		grpcAddr string
		timeout  string
	}

	rawPaymentClient struct {
		grpcAddr string
		timeout  string
	}

	rawRepo struct {
		postgresURL string
	}
)
