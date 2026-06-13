package config

import (
	"net"
	"net/url"
	"strings"
	"time"
)

func positiveDuration(value string) (time.Duration, error) {
	duration, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return 0, err
	}

	if duration <= 0 {
		return 0, errDurationMustBePositive
	}

	return duration, nil
}

func validateListen(addr string) error {
	_, port, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return err
	}

	if _, err := net.LookupPort("tcp", port); err != nil {
		return err
	}

	return nil
}

func validateGRPC(target string) error {
	if strings.Contains(target, "://") {
		return validateURL(target)
	}

	return validateAddr(target)
}

func validateAddr(addr string) error {
	host, port, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return err
	}

	if isBlank(host) {
		return errGRPCHostRequired
	}

	if _, err := net.LookupPort("tcp", port); err != nil {
		return err
	}

	return nil
}

func validateURL(target string) error {
	grpcTarget, err := url.Parse(strings.TrimSpace(target))
	if err != nil {
		return err
	}

	if grpcTarget.Scheme == "" {
		return errGRPCSchemeRequired
	}

	if !hasEndpoint(grpcTarget) {
		return errGRPCEndpointRequired
	}

	return nil
}

func postgresURL(rawURL string) (*url.URL, error) {
	postgresURL, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, err
	}

	if !postgresScheme(postgresURL) {
		return nil, errPostgresURLSchemeInvalid
	}

	if !hasDatabase(postgresURL) {
		return nil, errPostgresURLDatabaseRequired
	}

	return postgresURL, nil
}

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func hasEndpoint(target *url.URL) bool {
	return target.Opaque != "" || target.Host != "" || strings.Trim(target.Path, "/") != ""
}

func postgresScheme(postgresURL *url.URL) bool {
	return postgresURL.Scheme == "postgres" || postgresURL.Scheme == "postgresql"
}

func hasDatabase(postgresURL *url.URL) bool {
	return postgresURL.Host != "" && strings.Trim(postgresURL.Path, "/") != ""
}
