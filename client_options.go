package salt

import (
	"time"
)

type ClientOption func(options *Client)

func WithEndpoint(s string) ClientOption {
	return func(opt *Client) {
		opt.endpoint = s
	}
}

func WithInsecure() ClientOption {
	return func(opt *Client) {
		opt.skipVerify = true
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(opt *Client) {
		opt.timeout = timeout
	}
}

func WithUsername(username string) ClientOption {
	return func(opt *Client) {
		opt.username = username
	}
}

func WithPassword(password string) ClientOption {
	return func(opt *Client) {
		opt.password = password
	}
}

func WithAuthBackend(eauth string) ClientOption {
	return func(opt *Client) {
		opt.eauth = eauth
	}
}
