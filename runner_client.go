package salt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type ManageStatusReturn struct {
	UP   []string `json:"up"`
	Down []string `json:"down"`
}

func (c *Client) ManageStatus(ctx context.Context) (*ManageStatusReturn, error) {
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

func (c *Client) RunnerClient(ctx context.Context, fun string, arg []string, opts ...RunOption) (map[string]LocalClientReturn, error) {
	payload := commandRequest{
		Client:     RunnerClient,
		Target:     "*",
		Function:   fun,
		Arguments:  arg,
		TargetType: Glob,
		FullReturn: true,
	}

	for _, opt := range opts {
		opt(&payload)
	}

	data, err := c.post(ctx, "", payload)
	if err != nil {
		return nil, err
	}

	var resp localClientResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Return[0], nil
}

func (c *Client) RunnerClientAsync(ctx context.Context, tgt, fun string, arg []string, opts ...RunOption) (string, error) {
	payload := commandRequest{
		Client:     RunnerAsyncClient,
		Target:     tgt,
		Function:   fun,
		Arguments:  arg,
		TargetType: Glob,
	}
	for _, opt := range opts {
		opt(&payload)
	}

	data, err := c.post(ctx, "", payload)
	if err != nil {
		return "", err
	}

	if v := gjson.GetBytes(data, "status"); v.Value() != nil && v.Int() != 200 {
		return "", fmt.Errorf("failed to run async command, error: %s", gjson.GetBytes(data, "return").String())
	}

	jid := gjson.GetBytes(data, "return.0.jid")
	return jid.String(), nil
}
