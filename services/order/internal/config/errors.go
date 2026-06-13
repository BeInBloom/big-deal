package config

import "errors"

var (
	errDurationMustBePositive      = errors.New("must be positive")
	errGRPCHostRequired            = errors.New("host is required")
	errGRPCSchemeRequired          = errors.New("scheme is required")
	errGRPCEndpointRequired        = errors.New("endpoint is required")
	errPostgresURLSchemeInvalid    = errors.New("scheme must be postgres or postgresql")
	errPostgresURLDatabaseRequired = errors.New("host and database name are required")
)
