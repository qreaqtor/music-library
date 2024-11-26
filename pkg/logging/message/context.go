package logmsg

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const OperationID ContextKey = "operationID"

func ExtractOperationID(ctx context.Context) uuid.UUID {
	return ctx.Value(OperationID).(uuid.UUID)
}
