package salt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var _ Client = (*client)(nil)

type (
	Client interface {
		Login(ctx context.Context, username, password, eauth string) error
		Logout(ctx context.Context) error
		ListMinions(ctx context.Context) ([]Minion, error)
		GetMinion(ctx context.Context, mid string) (*Minion, error)
		LocalClient(ctx context.Context, tgt, fun string, arg []string, opts ...RunOption) (map[string]LocalClientReturn, error)
		LocalClientAsync(ctx context.Context, tgt, fun string, arg []string, opts ...RunOption) (jid string, err error)
		ListJobs(ctx context.Context) ([]Job, error)
		GetJob(ctx context.Context, jid string) (*Job, error)
		Hook(ctx context.Context, id string, payload interface{}) error
		Stats(ctx context.Context) (*stats, error)
		// Wheel Client Keys
		ListKeys(ctx context.Context) (*ListKeysReturn, error)
		GetKeyString(ctx context.Context, match string) (map[string]string, error)
		GetKeyFinger(ctx context.Context, match string) (map[string]string, error)
		AcceptKey(ctx context.Context, match string) ([]string, error)
		RejectedKey(ctx context.Context, match string) ([]string, error)
		DeleteKey(ctx context.Context, match string) ([]string, error)
	}
	client struct {
		httpClient *http.Client
		baseAPI    string
		token      string
	}
	clientOptions struct {
		skipVerify bool
		timeout    time.Duration
	}
	ClientOption func(options *clientOptions)
)

func WithInsecure() ClientOption {
	return func(options *clientOptions) {
		options.skipVerify = true
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *clientOptions) {
		options.timeout = timeout
	}
}

func NewClient(baseAPI string, opts ...ClientOption) *client {
	options := clientOptions{
		skipVerify: false,
		timeout:    60,
	}

	for _, o := range opts {
		o(&options)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	c := &client{
		baseAPI: baseAPI,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.skipVerify,
		},
	}
	c.httpClient = &http.Client{Transport: tr, Jar: jar, Timeout: options.timeout * time.Second}
	return c
}

func (c *client) doRequest(ctx context.Context, method, uri string, data interface{}) ([]byte, error) {
	url := strings.Join([]string{c.baseAPI, uri}, "/")

	var buf io.ReadWriter
	if data != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
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

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *client) get(ctx context.Context, uri string) ([]byte, error) {
	return c.doRequest(ctx, "GET", uri, nil)
}

func (c *client) post(ctx context.Context, uri string, data interface{}) ([]byte, error) {
	return c.doRequest(ctx, "POST", uri, data)
}
