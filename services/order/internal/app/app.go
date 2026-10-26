package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/BeInBloom/big-deal/services/order/internal/builder"
	"github.com/BeInBloom/big-deal/services/order/internal/config"
	"golang.org/x/sync/errgroup"
)

type App struct {
	server          *http.Server
	shutdownTimeout time.Duration
	closeResources  func() error
}

func New(cfg config.Config) (*App, error) {
	server, closeResources, err := builder.Build(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		server:          server,
		shutdownTimeout: cfg.HTTP().ShutdownTimeout(),
		closeResources:  closeResources,
	}, nil
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

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), a.shutdownTimeout,
		)
		defer cancel()

		//nolint:contextcheck // shutdown must outlive canceled run context
		err := a.server.Shutdown(shutdownCtx)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	})

	return errors.Join(g.Wait(), a.closeResources())
}
