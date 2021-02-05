package salt

import (
	"encoding/json"
	"fmt"
)

type KeysResponse struct {
	Return struct {
		Local           []string `json:"local"`
		Minions         []string `json:"minions"`
		MinionsPre      []string `json:"minions_pre"`
		MinionsRejected []string `json:"minions_rejected"`
		MinionsDenied   []string `json:"minions_denied"`
	} `json:"return"`
}

type KeyDetailResponse struct {
	Return struct {
		Minions map[string]string `json:"minions"`
	} `json:"return"`
}

func (c *Client) Keys() (*KeysResponse, error) {
	data, err := c.doRequest("GET", "keys", nil)
	if err != nil {
		return nil, err
	}

	var keys KeysResponse
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}

	return &keys, nil
}

func (c *Client) Key(mid string) (*KeyDetailResponse, error) {
	data, err := c.doRequest("GET", fmt.Sprintf("keys/%s", mid), nil)
	if err != nil {
		return nil, err
	}

	var keyDetail KeyDetailResponse

	if err := json.Unmarshal(data, &keyDetail); err != nil {
		return nil, err
	}

	return &keyDetail, nil
}
