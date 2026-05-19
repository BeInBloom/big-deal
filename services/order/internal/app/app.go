package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/BeInBloom/big-deal/services/order/internal/builder"
	"golang.org/x/sync/errgroup"
)

const (
	serverAddr      = ":8081"
	timeout         = time.Second * 5
	shutdownTimeout = time.Second * 2
)

type App struct {
	server *http.Server
}

func New() *App {
	orderHandlers, err := builder.Build()
	if err != nil {
		panic(err)
	}

	server := buildHttpServer(orderHandlers)

	return &App{
		server: server,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := buildNotifyContext(ctx)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	serverDone := make(chan struct{})

	g.Go(func() error {
		defer close(serverDone)

		log.Println("order service started")

		err := a.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
		case <-serverDone:
			return nil
		}

		log.Println("order service stopping")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := a.server.Shutdown(shutdownCtx)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	})

	return g.Wait()
}

func (a *App) Shutdown(ctx context.Context) error {
	err := a.server.Shutdown(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func buildNotifyContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}

func buildHttpServer(handlers http.Handler) *http.Server {
	return &http.Server{
		Addr:              serverAddr,
		Handler:           handlers,
		ReadHeaderTimeout: timeout,
	}

}
