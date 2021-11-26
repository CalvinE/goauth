package middleware

import (
	"net/http"

	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var RequestIDHeader = http.CanonicalHeaderKey("X-Request-Id")

// InitializeRequest initializes a request with a logger, request id, and trace
func InitializeRequest(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.Must(uuid.NewRandom()).String()
			}
			ctx = ctxpropagation.SetLoggerForContext(ctx, logger.With(
				zap.String("http_request_id", requestID),
			))
			ctx = ctxpropagation.SetRequestIDForContext(ctx, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
