package salt

import (
	"context"
	"encoding/json"
	"fmt"
)

type Job struct {
	Function   string      `json:"Function"`
	Arguments  []string    `json:"Arguments"`
	Target     interface{} `json:"Target"`
	TargetType string      `json:"Target-Type"`
	User       string      `json:"User"`
	StartTime  string      `json:"StartTime"`
}

type JobsResponse struct {
	Return []map[string]Job `json:"return"`
}

type JobResponse struct {
	Info []struct {
		*Job
		JID     string   `json:"jid"`
		Minions []string `json:"ListMinions"`
		Result  map[string]struct {
			Return  bool `json:"return"`
			Retcode int  `json:"retcode"`
			Success bool `json:"success"`
		} `json:"Result"`
	} `json:"info"`
	Return []map[string]bool `json:"return"`
}

func (c *client) ListJobs(ctx context.Context) (*JobsResponse, error) {
	data, err := c.doRequest(ctx,"GET", "jobs", nil)
	if err != nil {
		return nil, err
	}

	var jobs JobsResponse

	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, err
	}

	return &jobs, nil
}

func (c *client) GetJob(ctx context.Context, jid string) (*JobResponse, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("%s/%s", "jobs", jid), nil)
	if err != nil {
		return nil, err
	}

	var job JobResponse

	if err := json.Unmarshal(data, &job); err != nil {
		return nil, err
	}

	return &job, nil
}
