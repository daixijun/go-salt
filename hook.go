package salt

import (
	"context"
	"encoding/json"
	"errors"
)

type HookResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (c *Client) Hook(ctx context.Context, id string, payload interface{}) (*HookResponse, error) {
	data, err := c.doRequest(ctx, "POST", "hook/"+id, payload)
	if err != nil {
		return nil, err
	}

	var hook HookResponse

	if err := json.Unmarshal(data, &hook); err != nil {
		return nil, err
	}

	if !hook.Success {
		return nil, errors.New(hook.Message)
	}
	return &hook, nil

}
