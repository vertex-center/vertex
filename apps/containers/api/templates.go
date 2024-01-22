package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetTemplate(ctx context.Context, templateID string) (types.Template, error) {
	var template types.Template
	err := c.Request().
		Pathf("./templates/%s", templateID).
		ToJSON(&template).
		Fetch(ctx)
	return template, err
}
