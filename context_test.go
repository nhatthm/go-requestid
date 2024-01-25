package requestid_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.nhat.io/requestid"
)

func TestContext(t *testing.T) {
	t.Parallel()

	expected := uuid.New().String()

	ctx := requestid.ContextWithRequestID(context.Background(), expected)

	actual := requestid.FromContext(ctx)

	assert.Equal(t, expected, actual)
}
