package salt

import (
	"context"
	"encoding/json"
	"fmt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	EAuth    string `json:"eauth"`
}

type loginResponse struct {
	Return []struct {
		User   string   `json:"user"`
		Token  string   `json:"token"`
		EAuth  string   `json:"eauth"`
		Start  float64  `json:"start"`
		Expire float64  `json:"expire"`
		Perms  []string `json:"perms"`
	} `json:"return"`
}

// Login 认证
// https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#login
func (c *Client) Login(ctx context.Context) error {
	postData := loginRequest{
		Username: c.username,
		Password: c.password,
		EAuth:    c.eauth,
	}

	data, err := c.post(ctx, "login", postData)
	if err != nil {
		return err
	}

	var loginResp loginResponse
	if err := json.Unmarshal(data, &loginResp); err != nil {
		return err
	}

	c.token = loginResp.Return[0].Token

	return nil
}

// Logout 登出
// https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#logout
func (c *Client) Logout(ctx context.Context) error {
	_, err := c.post(ctx, "logout", nil)
	if err != nil {
		return fmt.Errorf("logout error: %w", err)
	}

	c.token = ""
	return nil
}
