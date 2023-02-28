package almlog

import (
	"bytes"
	"testing"

	"golang.org/x/exp/slog"
)

func UseTestLogger(t testing.TB) {
	opts := slog.HandlerOptions{
		Level:       Level,
		ReplaceAttr: removeTime,
	}
	Logger = slog.New(opts.NewTextHandler(tWriter{t}))
	slog.SetDefault(Logger)
}

type tWriter struct {
	t testing.TB
}

func (tw tWriter) Write(data []byte) (int, error) {
	tw.t.Log(string(bytes.TrimSpace(data)))
	return len(data), nil
}