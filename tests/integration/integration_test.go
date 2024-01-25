package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go.nhat.io/requestid"
)

func TestPropagator(t *testing.T) {
	t.Parallel()

	requestID := uuid.New().String()
	actual := ""

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actual = requestid.FromContext(r.Context())
	})

	srv := httptest.NewServer(otelhttp.NewHandler(h, "test",
		otelhttp.WithPropagators(requestid.Propagator{})),
	)

	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	req.Header.Set("x-request-id", requestID)

	resp, err := srv.Client().Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	assert.Equal(t, requestID, actual)
}

func TestHTTP(t *testing.T) {
	t.Parallel()

	actual := ""

	// Setup HTTP Client.
	c := &http.Client{
		Transport: requestid.DefaultTransport,
		Timeout:   time.Second,
	}

	// Spin up server 2.
	srv2 := httptest.NewServer(newRouter(func(w http.ResponseWriter, r *http.Request) {
		actual = requestid.FromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(srv2.Close)

	// Spin up server 1.
	srv1 := httptest.NewServer(newRouter(func(w http.ResponseWriter, r *http.Request) {
		nextReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, srv2.URL, nil)
		require.NoError(t, err)

		resp, err := c.Do(nextReq)
		require.NoError(t, err)

		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(srv1.Close)

	// Make request for testing.
	requestID := uuid.New().String()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv1.URL, nil)
	require.NoError(t, err)

	req.Header.Set("x-request-id", requestID)

	resp, err := c.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	// Assertions.
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, requestID, actual)
}

func newRouter(h func(w http.ResponseWriter, r *http.Request)) *chi.Mux {
	r := chi.NewRouter()

	r.Use(requestid.HandlerMiddleware())
	r.Get("/", h)

	return r
}
