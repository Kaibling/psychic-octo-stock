package apimiddleware

import (
	"context"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/lucsky/cuid"
)

func Response(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestID string
		clientRequestID := r.Header.Get("X-REQUEST-ID")
		if clientRequestID != "" {
			requestID = clientRequestID
		} else {
			requestID = cuid.Slug()
		}
		response := transmission.NewResponse(w, r, requestID)
		ctx := context.WithValue(r.Context(), "responseObject", response)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
