package salt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Client interface {
	Login(ctx context.Context, username, password, eauth string) error
	Logout(ctx context.Context) error
	ListMinions(ctx context.Context) (*MinionResponse, error)
	GetMinion(ctx context.Context, mid string) (*MinionResponse, error)
	ListKeys(ctx context.Context) (*KeysResponse, error)
	GetKey(ctx context.Context, mid string) (*KeyDetailResponse, error)
	ListJobs(ctx context.Context)(*JobsResponse, error)
	GetJob(ctx context.Context, jid string) (*JobResponse, error)
	Run(ctx context.Context, payload *RunRequest)(*RunResponse, error)
	Hook(ctx context.Context, id string, payload interface{}) (*HookResponse, error)
	Stats(ctx context.Context) (*StatsResponse, error)
}

type client struct {
	httpClient   *http.Client
	BaseAPI      string
	ExternalAuth *ExternalAuth
	// Token      string
}

func NewClient(baseAPI string, skipVerify bool) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipVerify,
		},
	}
	cookieJar, _ := cookiejar.New(nil)
	return &client{
		httpClient: &http.Client{Transport: tr, Jar: cookieJar},
		BaseAPI:    baseAPI,
	}
}

func (c *client) doRequest(ctx context.Context, method, uri string, data interface{}) ([]byte, error) {
	// url := fmt.Sprintf("%s/%s", testClient.BaseAPI, uri)
	url := strings.Join([]string{c.BaseAPI, uri}, "/")

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
	//
	// if c.Token != "" {
	// 	req.Header.Set("X-Auth-Token", testClient.Token)
	// }

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
