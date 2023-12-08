package api

import (
	"context"
)

func (c *Client) Reboot(ctx context.Context) error {
	return c.Request().
		Pathf("./hardware/reboot").
		Post().
		Fetch(ctx)
}

func (c *KernelClient) Reboot(ctx context.Context) error {
	return c.Request().
		Pathf("./hardware/reboot").
		Post().
		Fetch(ctx)
}
