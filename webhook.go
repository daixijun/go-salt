package salt

import (
	"context"
	"encoding/json"
	"fmt"
)

type hookResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (c *Client) Hook(ctx context.Context, id string, payload interface{}) error {
	data, err := c.post(ctx, "hook/"+id, payload)
	if err != nil {
		return fmt.Errorf("failed to post hook: %w", err)
	}

	var resp hookResponse

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal hook response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("unexpected hook response: %s, status: %d", resp.Message, resp.Status)
	}
	return nil
}
