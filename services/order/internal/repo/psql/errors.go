package psql

import (
	"errors"
	"fmt"

	"github.com/BeInBloom/big-deal/services/order/internal/repo/psql/gen/dberrors"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	errOrderPartsNotLoaded  = errors.New("order parts are not loaded")
	errUnknownOrderStatus   = errors.New("unknown order status")
	errUnknownPaymentMethod = errors.New("unknown payment method")
	errInvalidPartPrice     = errors.New("invalid part price")
)

func intoCreateOrderError(err error) error {
	// The arguments are reversed on purpose: the bob-generated constraint error
	// implements Is(*pgconn.PgError), not the other way around.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && errors.Is(dberrors.OrderErrors.ErrUniqueOrdersPkey, pgErr) {
		return fmt.Errorf("%w: %w", services.ErrOrderAlreadyExists, err)
	}

	return err
}
