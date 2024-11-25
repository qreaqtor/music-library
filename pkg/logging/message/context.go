package logmsg

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

func ExtractLogID(ctx context.Context) uuid.UUID {
	return ctx.Value(ContextKey("logID")).(uuid.UUID)
}
