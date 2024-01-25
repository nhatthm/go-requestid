package requestid

import (
	"context"

	"github.com/bool64/ctxd"
)

type contextKey struct{}

// ContextWithRequestID injects the request ID into the context.
func ContextWithRequestID(ctx context.Context, id string) context.Context {
	ctx = context.WithValue(ctx, contextKey{}, id)
	ctx = ctxd.AddFields(ctx, "request_id", id)

	return ctx
}

// FromContext extracts the request ID from the context.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string) //nolint: errcheck

	return id
}
