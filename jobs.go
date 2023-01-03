package salt

import (
	"context"
	"encoding/json"
	"fmt"
)

type Job struct {
	JID       string        `json:"jid,omitempty"`
	Function  string        `json:"Function,omitempty"`
	Arguments []interface{} `json:"Arguments,omitempty"`
	// KWArguments map[string]interface{} `json:"KWArguments,omitempty"`
	Target     interface{}          `json:"Target,omitempty"`
	TargetType TargetType           `json:"Target-Type,omitempty"`
	StartTime  saltTime             `json:"StartTime"`
	User       string               `json:"User"`
	Minions    []string             `json:"Minions,omitempty"`
	Result     map[string]jobResult `json:"Result,omitempty"`
	Error      string               `json:"Error,omitempty"`
}

type jobResult struct {
	Return  interface{} `json:"return"`
	Retcode int         `json:"retcode"`
	Success bool        `json:"success"`
	Out     string      `json:"out,omitempty"`
}

type jobDetailResponse struct {
	Info   []Job                    `json:"info"`
	Return []map[string]interface{} `json:"return"`
}

type jobListResponse struct {
	Return []map[string]Job `json:"return"`
}

// ListJobs 获取任务列表
// https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#jobs
func (c *Client) ListJobs(ctx context.Context) ([]Job, error) {
	data, err := c.get(ctx, "jobs")
	if err != nil {
		return nil, err
	}

	var resp jobListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	jobs := make([]Job, len(resp.Return[0]))
	for jid, job := range resp.Return[0] {
		job.JID = jid
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// LookupJID 获取任务详情
// https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--jobs-(jid)
func (c *Client) LookupJID(ctx context.Context, jid string) (*Job, error) {
	data, err := c.get(ctx, fmt.Sprintf("%s/%s", "jobs", jid))
	if err != nil {
		return nil, err
	}

	var resp jobDetailResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	info := resp.Info[0]
	if info.Error != "" {
		return nil, fmt.Errorf("cannot receive job result, error message: %s", info.Error)
	}
	return &info, nil
}
