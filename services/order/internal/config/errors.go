package config

import "errors"

var (
	errDurationMustBePositive      = errors.New("must be positive")
	errGRPCHostRequired            = errors.New("host is required")
	errGRPCSchemeRequired          = errors.New("scheme is required")
	errGRPCEndpointRequired        = errors.New("endpoint is required")
	errRepoKindRequired            = errors.New("is required")
	errRepoKindUnknown             = errors.New("must be hash_map or postgres")
	errPostgresURLRequired         = errors.New("is required for postgres repo")
	errPostgresURLSchemeInvalid    = errors.New("scheme must be postgres or postgresql")
	errPostgresURLDatabaseRequired = errors.New("host and database name are required")
)
