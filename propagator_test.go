package requestid_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"

	"go.nhat.io/requestid"
)

func TestPropagator_Inject(t *testing.T) {
	t.Parallel()

	requestID := uuid.New().String()

	testCases := []struct {
		scenario string
		context  context.Context
		expected map[string]string
	}{
		{
			scenario: "no request id",
			context:  context.Background(),
			expected: map[string]string{},
		},
		{
			scenario: "has request id",
			context:  requestid.ContextWithRequestID(context.Background(), requestID),
			expected: map[string]string{
				"x-request-id": requestID,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := map[string]string{}

			requestid.Propagator{}.Inject(tc.context, propagation.MapCarrier(actual))

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPropagator_Extract(t *testing.T) {
	t.Parallel()

	requestID := uuid.New().String()

	testCases := []struct {
		scenario string
		carrier  map[string]string
		expected string
	}{
		{
			scenario: "no request id",
			carrier:  map[string]string{},
		},
		{
			scenario: "has request id",
			carrier: map[string]string{
				"x-request-id": requestID,
			},
			expected: requestID,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			ctx := requestid.Propagator{}.Extract(context.Background(), propagation.MapCarrier(tc.carrier))
			actual := requestid.FromContext(ctx)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPropagator_Fields(t *testing.T) {
	t.Parallel()

	actual := requestid.Propagator{}.Fields()
	expected := []string{"x-request-id"}

	assert.Equal(t, expected, actual)
}
