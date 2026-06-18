package psql

import "errors"

var (
	errOrderPartsNotLoaded  = errors.New("order parts are not loaded")
	errUnknownOrderStatus   = errors.New("unknown order status")
	errUnknownPaymentMethod = errors.New("unknown payment method")
	errInvalidPartPrice     = errors.New("invalid part price")
)
