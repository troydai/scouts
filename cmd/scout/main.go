package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/troydai/scouts/internal/httpserver"
)

func main() {
	support := fx.Options(
		fx.Provide(zap.NewDevelopment),
	)

	core := fx.Options(
		httpserver.Module,
	)

	fx.New(support, core).Run()
}
