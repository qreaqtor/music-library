package logmsg

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type LogMsg struct {
	id     uuid.UUID
	URL    string
	Method string
	Text   string
	Status int
}

// Возвращает структуру, которая пишет логи с помощью logger.
// Остальные поля - информация, которая будет выводиться.
func NewLogMsg(ctx context.Context, url, method string) *LogMsg {
	return &LogMsg{
		URL:    url,
		Method: method,
		id:ExtractLogID(ctx),
	}
}

func (msg *LogMsg) With(text string, status int) *LogMsg {
	msg.Text = text
	msg.Status = status
	return msg
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
	}
}
