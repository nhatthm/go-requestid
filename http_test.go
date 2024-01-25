package requestid_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/requestid"
)

func TestHandlerMiddleware_NoRequestID(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		actual := requestid.FromContext(req.Context())

		assert.NotEmpty(t, actual)
	})

	srv := httptest.NewServer(requestid.HandlerMiddleware()(h))

	t.Cleanup(srv.Close)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
}

func TestHandlerMiddleware_HasRequestID(t *testing.T) {
	t.Parallel()

	requestID := uuid.New().String()

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		actual := requestid.FromContext(req.Context())

		assert.Equal(t, requestID, actual)
	})

	srv := httptest.NewServer(requestid.HandlerMiddleware()(h))

	t.Cleanup(srv.Close)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	req.Header.Set("x-request-id", requestID)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
}

func TestRoundTripperMiddleware(t *testing.T) {
	t.Parallel()

	requestID := uuid.New().String()

	testCases := []struct {
		scenario string
		context  context.Context
		expected string
	}{
		{
			scenario: "no request id",
			context:  context.Background(),
			expected: "",
		},
		{
			scenario: "has request id",
			context:  requestid.ContextWithRequestID(context.Background(), requestID),
			expected: requestID,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var actual string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				actual = req.Header.Get("x-request-id")
			}))

			t.Cleanup(srv.Close)

			c := http.Client{Transport: requestid.RoundTripperMiddleware()(http.DefaultTransport)}

			req, err := http.NewRequestWithContext(tc.context, http.MethodGet, srv.URL, nil)
			require.NoError(t, err)

			resp, err := c.Do(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())

			assert.Equal(t, tc.expected, actual)
		})
	}
}
