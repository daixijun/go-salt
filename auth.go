package salt

import (
	"encoding/json"
)

type ExternalAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	EAuth    string `json:"eauth"`
}

type LoginResponse struct {
	Return []struct {
		Token  string   `json:"Token"`
		Start  float64  `json:"start"`
		Expire float64  `json:"expire"`
		User   string   `json:"user"`
		EAuth  string   `json:"eauth"`
		Perms  []string `json:"perms"`
	} `json:"return"`
}

// Login 认证
func (c *Client) Login(username, password, eauth string) error {
	externalAuth := ExternalAuth{
		Username: username,
		Password: password,
		EAuth:    eauth,
	}

	data, err := c.doRequest("POST", "login", externalAuth)
	if err != nil {
		return err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(data, &loginResp); err != nil {
		return err
	}

	c.ExternalAuth = &externalAuth

	return nil
}

func (c *Client) Logout() error {
	_, err := c.doRequest("POST", "logout", nil)
	return err
}
