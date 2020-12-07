package salt

import (
	"context"
	"encoding/json"
)

type RunRequest struct {
	*ExternalAuth
	Client   string            `json:"client"`
	Target   interface{}       `json:"tgt"`
	Function string            `json:"fun"`
	JID      string            `json:"jid,omitempty"` // runner
	Arg      []string          `json:"arg,omitempty"`
	KwArg    map[string]string `json:"kwarg,omitempty"`
}

type RunResponse struct {
	Return []map[string]interface{} `json:"return"`
}

func (c *Client) Run(ctx context.Context, payload *RunRequest) (*RunResponse, error) {
	payload.ExternalAuth = c.ExternalAuth

	data, err := c.doRequest(ctx, "POST", "/run", payload)
	if err != nil {
		return nil, err
	}

	var ret RunResponse

	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
