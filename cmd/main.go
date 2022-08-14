package main

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	support := fx.Options(
		fx.Provide(zap.NewDevelopment),
	)

	core := fx.Options(
		fx.Invoke(run),
	)

	app := fx.New(support, core)

	app.Run()
}

func run(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("hey, i'm scout.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("bye")
			return nil
		},
	})
}
