package salt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
	token      string
	skipVerify bool
	timeout    time.Duration
	username   string
	password   string
	eauth      string
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		endpoint:   "https://localhost:8000",
		skipVerify: false,
		timeout:    60,
		username:   "salt",
		password:   "salt",
		eauth:      "pam",
	}

	for _, o := range opts {
		o(c)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.skipVerify,
		},
	}
	c.httpClient = &http.Client{Transport: tr, Jar: jar, Timeout: c.timeout * time.Second}
	return c
}

func (c *Client) doRequest(ctx context.Context, method, uri string, data interface{}) ([]byte, error) {
	url := strings.Join([]string{c.endpoint, uri}, "/")

	var buf bytes.Buffer
	if data != nil {
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.token != "" {
		req.Header.Set("X-Auth-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) get(ctx context.Context, uri string) ([]byte, error) {
	return c.doRequest(ctx, "GET", uri, nil)
}

func (c *Client) post(ctx context.Context, uri string, data interface{}) ([]byte, error) {
	return c.doRequest(ctx, "POST", uri, data)
}
