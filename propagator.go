package requestid

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

const fieldName = "x-request-id"

var _ propagation.TextMapPropagator = (*Propagator)(nil)

// Propagator is a request ID propagator.
type Propagator struct{}

// Inject injects the request ID into the header.
func (p Propagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	if id := FromContext(ctx); id != "" {
		carrier.Set(fieldName, id)
	}
}

// Extract extracts the request ID from the header.
func (p Propagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	if id := carrier.Get(fieldName); id != "" {
		ctx = ContextWithRequestID(ctx, id)
	}

	return ctx
}

// Fields returns the keys whose values are set with ContextWithRequestID.
func (p Propagator) Fields() []string {
	return []string{fieldName}
}
