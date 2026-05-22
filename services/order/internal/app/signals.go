package app

import (
	"context"
	"os/signal"
	"syscall"
)

func buildNotifyContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}
