package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	_port = 8080
)

var Module = fx.Options(
	fx.Provide(ProvideMux),
	fx.Invoke(Register),
)

func Register(lc fx.Lifecycle, logger *zap.Logger, mux http.Handler) {
	serverExitCh := make(chan error)
	var server *http.Server

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", _port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("fail to listen to TCP: %w", err)
			}

			server := &http.Server{
				Addr:    addr,
				Handler: mux,
			}

			go func() {
				defer close(serverExitCh)
				serverExitCh <- server.Serve(lis)
				logger.Info("HTTP server exited")
			}()

			logger.Info("Start listening", zap.String("address", addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if server == nil {
				return nil
			}

			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("fail to stop server: %w", err)
			}

			select {
			case <-ctx.Done():
				logger.Warn("Terminate before the service is stopped due to timeout.")
			case err := <-serverExitCh:
				logger.Info("Stopped listenning", zap.Error(err))
			}

			return nil
		},
	})
}
