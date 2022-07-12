package salt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type ListKeysReturn struct {
	Local           []string `json:"local"`
	Minions         []string `json:"minions"`
	MinionsPre      []string `json:"minions_pre"`
	MinionsRejected []string `json:"minions_rejected"`
	MinionsDenied   []string `json:"minions_denied"`
}

type listKeysResponse struct {
	Return ListKeysReturn `json:"return"`
}

func (c *client) ListKeys(ctx context.Context) (*ListKeysReturn, error) {
	data, err := c.get(ctx, "keys")
	if err != nil {
		return nil, err
	}

	var resp listKeysResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp.Return, nil
}

func (c *client) GetKeyString(ctx context.Context, match string) (map[string]string, error) {
	req := commandRequest{
		Client:   WheelClient,
		Function: "key.print",
		Match:    match,
	}
	data, err := c.post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	if v := gjson.GetBytes(data, "return.0.data.success"); !v.Bool() {
		return nil, fmt.Errorf("failed to get key for %s, error: %s", match, gjson.GetBytes(data, "return.0.data.return").String())
	}

	value := gjson.GetBytes(data, "return.0.data.return.minions")

	minions := make(map[string]string)
	for k, v := range value.Map() {
		minions[k] = v.String()
	}
	return minions, nil
}

func (c *client) GetKeyFinger(ctx context.Context, match string) (map[string]string, error) {
	req := commandRequest{
		Client:   WheelClient,
		Function: "key.finger",
		Match:    match,
	}
	data, err := c.post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	if v := gjson.GetBytes(data, "return.0.data.success"); !v.Bool() {
		return nil, fmt.Errorf("failed to get key finger for %s, error: %s", match, gjson.GetBytes(data, "return.0.data.return").String())
	}

	value := gjson.GetBytes(data, "return.0.data.return.minions."+match)
	minions := make(map[string]string)
	for k, v := range value.Map() {
		minions[k] = v.String()
	}
	return minions, nil
}

func (c *client) AcceptKey(ctx context.Context, mid string) ([]string, error) {
	req := commandRequest{
		Client:   WheelClient,
		Function: "key.accept",
		Match:    mid,
	}
	data, err := c.post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	if v := gjson.GetBytes(data, "return.0.data.success"); !v.Bool() {
		return nil, fmt.Errorf("failed to accept key for %s, error: %s", mid, gjson.GetBytes(data, "return.0.data.return").String())
	}

	value := gjson.GetBytes(data, "return.0.data.return.minions")

	var minions []string
	for _, v := range value.Array() {
		minions = append(minions, v.String())
	}
	return minions, nil
}

func (c *client) RejectedKey(ctx context.Context, mid string) ([]string, error) {
	req := commandRequest{
		Client:   WheelClient,
		Function: "key.reject",
		Match:    mid,
	}
	data, err := c.post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	if v := gjson.GetBytes(data, "return.0.data.success"); !v.Bool() {
		return nil, fmt.Errorf("failed to reject key for %s, error: %s", mid, gjson.GetBytes(data, "return.0.data.return").String())
	}

	value := gjson.GetBytes(data, "return.0.data.return.minions")

	var minions []string
	for _, v := range value.Array() {
		minions = append(minions, v.String())
	}
	return minions, nil
}

func (c *client) DeleteKey(ctx context.Context, match string) ([]string, error) {
	req := commandRequest{
		Client:   WheelClient,
		Function: "key.delete",
		Match:    match,
	}
	data, err := c.post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	if v := gjson.GetBytes(data, "return.0.data.success"); !v.Bool() {
		return nil, fmt.Errorf("failed to delete key for %s, error: %s", match, gjson.GetBytes(data, "return.0.data.return").String())
	}
	value := gjson.GetBytes(data, "return.0.data.return.minions")

	var minions []string
	for _, v := range value.Array() {
		minions = append(minions, v.String())
	}
	return minions, nil
}
