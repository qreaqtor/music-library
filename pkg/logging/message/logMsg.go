package logmsg

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type LogMsg struct {
	opareation uuid.UUID
	URL        string
	Method     string
	Text       string
	Status     int
}

// Возвращает структуру, которая пишет логи с помощью logger.
// Остальные поля - информация, которая будет выводиться.
func NewLogMsg(ctx context.Context, url, method string) *LogMsg {
	return &LogMsg{
		URL:        url,
		Method:     method,
		opareation: ExtractOperationID(ctx),
	}
}

func (msg *LogMsg) With(text string, status int) *LogMsg {
	return &LogMsg{
		Text:       text,
		Status:     status,
		opareation: msg.opareation,
		URL:        msg.URL,
		Method:     msg.Method,
	}
}

func (msg *LogMsg) Info() {
	slog.Info(msg.Text, getArgs(msg)...)
}

func (msg *LogMsg) Error() {
	slog.Error(msg.Text, getArgs(msg)...)
}

func getArgs(msg *LogMsg) []any {
	return []any{
		"status", msg.Status,
		"url", msg.URL,
		"method", msg.Method,
		"operation", msg.opareation,
	}
}
