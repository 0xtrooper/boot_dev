package middleware

import "log/slog"

type Middleware struct {
	Metrics metrics
}

func NewMiddleware(l *slog.Logger) (Middleware) {

	return Middleware{
		Metrics: metrics{
			logger: l.With("middleware", "metrics"),
			fileserverHits: 0,
		},
	}
}
