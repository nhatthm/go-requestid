# RequestID

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-requestid)](https://github.com/nhatthm/go-requestid/releases/latest)
[![Build Status](https://github.com/nhatthm/go-requestid/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-requestid/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-requestid/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/requestid)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/requestid)](https://goreportcard.com/report/go.nhat.io/requestid)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/requestid)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

A library to propagate Request ID across the context.

## Prerequisites

- `Go >= 1.22`

## Install

```bash
go get go.nhat.io/requestid
```

## Usage

### HTTP Server

You can use `requestid.NewHandler()` to wrap your HTTP handler or `requestid.HandlerMiddleware()` if you use `chi` router.

### HTTP Client

You can use `requestid.DefaultTransport` or `requestid.NewRoundTripper()` to wrap your HTTP transport. If your client factory supports middleware, you can use `requestid.RoundTripperMiddleware()`.

### OpenTelemetry Propagation

You can use `requestid.Propagator` to propagate Request ID across the context. See [the doc](https://opentelemetry.io/docs/specs/otel/context/api-propagators/) for more details.

## Examples

Example 1: HTTP server and client middleware.

```go
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

	"go.nhat.io/requestid"
)

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
```

Example 2: OpenTelemetry propagation.

```go
package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
