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

type Client struct {
	Ctx          context.Context
	httpClient   *http.Client
	BaseAPI      string
	ExternalAuth *ExternalAuth
	// Token      string
}

func NewClient(ctx context.Context, baseAPI string, skipVerify bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipVerify,
		},
	}
	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		Ctx:        ctx,
		httpClient: &http.Client{Transport: tr, Jar: cookieJar},
		BaseAPI:    baseAPI,
	}
}

func (c *Client) doRequest(method, uri string, data interface{}) ([]byte, error) {
	// url := fmt.Sprintf("%s/%s", c.BaseAPI, uri)
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

	req, err := http.NewRequestWithContext(c.Ctx, method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	//
	// if c.Token != "" {
	// 	req.Header.Set("X-Auth-Token", c.Token)
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
