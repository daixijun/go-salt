package salt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type LocalClientReturn struct {
	Jid     string      `json:"jid"`
	Ret     interface{} `json:"ret"`
	Retcode int         `json:"retcode"`
}
type localClientResponse struct {
	Return []map[string]LocalClientReturn `json:"return"`
}

func (c *Client) LocalClient(ctx context.Context, fun string, arg []string, opts ...RunOption) (map[string]LocalClientReturn, error) {
	payload := commandRequest{
		Client:     LocalClient,
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

func (c *Client) LocalClientAsync(ctx context.Context, tgt, fun string, arg []string, opts ...RunOption) (string, error) {
	payload := commandRequest{
		Client:     LocalAsyncClient,
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
