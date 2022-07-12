package salt

import (
	"context"
	"encoding/json"
)

type ManageStatusReturn struct {
	UP   []string `json:"up"`
	Down []string `json:"down"`
}

func (c *client) ManageStatus(ctx context.Context) (*ManageStatusReturn, error) {
	req := commandRequest{
		Client:   RunnerClient,
		Function: "manage.status",
	}
	data, err := c.post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	var ret ManageStatusReturn
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
