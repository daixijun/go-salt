package salt

import (
	"context"
	"encoding/json"
)

type HookResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (c *client) Hook(ctx context.Context, id string, payload interface{}) (*HookResponse, error) {
	data, err := c.doRequest(ctx, "POST", "hook/"+id, payload)
	if err != nil {
		return nil, err
	}

	var hook *HookResponse

	err = json.Unmarshal(data, &hook)
	return hook, err
}
