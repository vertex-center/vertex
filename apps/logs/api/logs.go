package api

import (
	"golang.org/x/net/context"
)

func (c *Client) PushLogs(ctx context.Context, content string) error {
	return c.Request().
		Pathf("./logs/push").
		BodyJSON(map[string]interface{}{
			"content": content,
		}).
		Post().
		Fetch(ctx)
}
