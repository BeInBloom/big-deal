package config

import (
	"net/url"
	"time"
)

func (c config) HTTP() HTTP {
	return c.http
}

func (c config) Inventory() InventoryClient {
	return c.inventory
}

func (c config) Payment() PaymentClient {
	return c.payment
}

func (c config) Repo() Repo {
	return c.repo
}

func (c config) isConfig() {}

func (h httpConfig) Addr() string {
	return h.addr
}

func (h httpConfig) ReadHeaderTimeout() time.Duration {
	return h.readHeaderTimeout
}

func (h httpConfig) ShutdownTimeout() time.Duration {
	return h.shutdownTimeout
}

func (h httpConfig) isHTTP() {}

func (c inventoryClient) GRPCAddr() string {
	return c.grpcAddr
}

func (c inventoryClient) Timeout() time.Duration {
	return c.timeout
}

func (c inventoryClient) isInventoryClient() {}

func (c paymentClient) GRPCAddr() string {
	return c.grpcAddr
}

func (c paymentClient) Timeout() time.Duration {
	return c.timeout
}

func (c paymentClient) isPaymentClient() {}

func (r repoConfig) IsInMemory() bool {
	return r.postgresURL == nil
}

func (r repoConfig) PostgresURL() (url.URL, bool) {
	if r.postgresURL == nil {
		return url.URL{}, false
	}

	return *r.postgresURL, true
}

func (r repoConfig) isRepo() {}
