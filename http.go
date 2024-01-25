package requestid

import (
	"net/http"

	"github.com/google/uuid"
)

// DefaultTransport is a wrapper around http.DefaultTransport that injects the request ID into the header.
var DefaultTransport = NewRoundTripper(http.DefaultTransport)

// NewHandler wraps a http.Handler to extract the requestID from the header and inject it into the context.
func NewHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := req.Header.Get(fieldName)
		if id == "" {
			id = uuid.New().String()

			req.Header.Set(fieldName, id)
		}

		ctx := ContextWithRequestID(req.Context(), id)

		h.ServeHTTP(w, req.WithContext(ctx))
	})
}

// HandlerMiddleware returns a middleware that extracts the requestID from the header and injects it into the context.
func HandlerMiddleware() func(http.Handler) http.Handler {
	return NewHandler
}

var _ http.RoundTripper = (*roundTripper)(nil)

type roundTripper struct {
	rt http.RoundTripper
}

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if id := FromContext(req.Context()); id != "" {
		req.Header.Set(fieldName, id)
	}

	return r.rt.RoundTrip(req) //nolint: wrapcheck
}

// NewRoundTripper wraps a http.RoundTripper to inject the request ID in the context into the header.
func NewRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return roundTripper{
		rt: rt,
	}
}

// RoundTripperMiddleware returns a middleware that injects the request ID in the context into the header.
func RoundTripperMiddleware() func(http.RoundTripper) http.RoundTripper {
	return NewRoundTripper
}
