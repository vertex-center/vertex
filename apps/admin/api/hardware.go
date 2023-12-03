package api

import (
	"context"

	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) Reboot(ctx context.Context) error {
	var apiError api.Error
	err := c.Request().
		Pathf("./hardware/reboot").
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return err
}

func (c *KernelClient) Reboot(ctx context.Context) error {
	var apiError api.Error
	err := c.Request().
		Pathf("./hardware/reboot").
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return err
}
