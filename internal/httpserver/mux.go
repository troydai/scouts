package httpserver

import (
	"net/http"

	"go.uber.org/zap"
)

func ProvideMux(logger *zap.Logger) http.Handler {
	svc := service{logger: logger}
	mux := http.NewServeMux()
	mux.HandleFunc("/_health", svc.Health)

	return mux
}

type service struct {
	logger *zap.Logger
}

func (s *service) Health(w http.ResponseWriter, _ *http.Request) {
	s.logger.Debug("inbound request", zap.String("path", "/_health"))
	w.WriteHeader(http.StatusOK)
}
