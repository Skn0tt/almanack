package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/carlmjohnson/slackhook"
	"golang.org/x/exp/slog"
)

type Client struct {
	c *slackhook.Client
}

type (
	Message    = slackhook.Message
	Attachment = slackhook.Attachment
	Field      = slackhook.Field
)

func New(hookURL string) Client {
	return Client{slackhook.New(hookURL, nil)}
}

func (sc Client) Post(ctx context.Context, msg Message) error {
	l := slog.FromContext(ctx)
	if sc.c == nil {
		l.Info("slack.Post: mocking; debug output")
		b, _ := json.MarshalIndent(&msg, "", "  ")
		fmt.Fprintf(os.Stderr, "\n%s\n", b)
		return nil
	}
	l.Info("slack.Post")
	return sc.c.PostCtx(ctx, msg)
}
