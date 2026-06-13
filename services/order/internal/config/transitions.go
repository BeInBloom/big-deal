package config

import (
	"fmt"
	"strings"
	"time"
)

func (r raw) tryIntoConfig() (Config, error) {
	http, err := r.http.tryIntoHTTP()
	if err != nil {
		return nil, err
	}

	inventory, err := r.inventory.tryIntoInventoryClient()
	if err != nil {
		return nil, err
	}

	payment, err := r.payment.tryIntoPaymentClient()
	if err != nil {
		return nil, err
	}

	repo, err := r.repo.tryIntoRepo()
	if err != nil {
		return nil, err
	}

	return config{
		http:      http,
		inventory: inventory,
		payment:   payment,
		repo:      repo,
	}, nil
}

func (r rawHTTP) tryIntoHTTP() (httpConfig, error) {
	addr := strings.TrimSpace(r.addr)
	if addr == "" {
		addr = defaultHTTPAddr
	}

	if err := validateListen(addr); err != nil {
		return httpConfig{}, fmt.Errorf("%s: %w", envHTTPAddr, err)
	}

	readHeaderTimeout, err := r.readHeaderDuration()
	if err != nil {
		return httpConfig{}, fmt.Errorf("%s: %w", envHTTPReadHeaderTimeout, err)
	}

	shutdownTimeout, err := r.shutdownDuration()
	if err != nil {
		return httpConfig{}, fmt.Errorf("%s: %w", envHTTPShutdownTimeout, err)
	}

	return httpConfig{
		addr:              addr,
		readHeaderTimeout: readHeaderTimeout,
		shutdownTimeout:   shutdownTimeout,
	}, nil
}

func (r rawInventoryClient) tryIntoInventoryClient() (inventoryClient, error) {
	grpcAddr := strings.TrimSpace(r.grpcAddr)
	if grpcAddr == "" {
		grpcAddr = defaultInventoryAddr
	}

	if err := validateGRPC(grpcAddr); err != nil {
		return inventoryClient{}, fmt.Errorf("%s: %w", envInventoryGRPCAddr, err)
	}

	timeout, err := r.timeoutDuration()
	if err != nil {
		return inventoryClient{}, fmt.Errorf("%s: %w", envInventoryTimeout, err)
	}

	return inventoryClient{
		grpcAddr: grpcAddr,
		timeout:  timeout,
	}, nil
}

func (r rawPaymentClient) tryIntoPaymentClient() (paymentClient, error) {
	grpcAddr := strings.TrimSpace(r.grpcAddr)
	if grpcAddr == "" {
		grpcAddr = defaultPaymentAddr
	}

	if err := validateGRPC(grpcAddr); err != nil {
		return paymentClient{}, fmt.Errorf("%s: %w", envPaymentGRPCAddr, err)
	}

	timeout, err := r.timeoutDuration()
	if err != nil {
		return paymentClient{}, fmt.Errorf("%s: %w", envPaymentTimeout, err)
	}

	return paymentClient{
		grpcAddr: grpcAddr,
		timeout:  timeout,
	}, nil
}

func (r rawHTTP) readHeaderDuration() (time.Duration, error) {
	if isBlank(r.readHeaderTimeout) {
		return defaultReadHeaderTimeout, nil
	}

	return positiveDuration(r.readHeaderTimeout)
}

func (r rawHTTP) shutdownDuration() (time.Duration, error) {
	if isBlank(r.shutdownTimeout) {
		return defaultShutdownTimeout, nil
	}

	return positiveDuration(r.shutdownTimeout)
}

func (r rawInventoryClient) timeoutDuration() (time.Duration, error) {
	if isBlank(r.timeout) {
		return defaultClientTimeout, nil
	}

	return positiveDuration(r.timeout)
}

func (r rawPaymentClient) timeoutDuration() (time.Duration, error) {
	if isBlank(r.timeout) {
		return defaultClientTimeout, nil
	}

	return positiveDuration(r.timeout)
}

func (r rawRepo) tryIntoRepo() (repoConfig, error) {
	if isBlank(r.postgresURL) {
		return repoConfig{}, nil
	}

	pgURL, err := postgresURL(r.postgresURL)
	if err != nil {
		return repoConfig{}, fmt.Errorf("%s: %w", envRepoPostgresURL, err)
	}

	return repoConfig{
		postgresURL: pgURL,
	}, nil
}
