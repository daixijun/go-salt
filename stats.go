package salt

import (
	"encoding/json"
	"strings"
)

const (
	CP_APPLICATIONS      = "CherryPy Applications"
	CP_HTTPSERVER_PREFIX = "Cheroot HTTPServer"
)

type TMPStats map[string]interface{}

type StatsResponse struct {
	Applications Applications `json:"CherryPy Applications"`
	HTTPServer   HTTPServer   `json:"Cheroot HTTPServer"`
}

type Applications struct {
	Enabled             bool    `json:"Enabled"`
	BytesReadRequest    float64 `json:"Bytes Read/Request"`
	BytesReadSecond     float64 `json:"Bytes Read/Second"`
	BytesWrittenRequest float64 `json:"Bytes Written/Request"`
	BytesWrittenSecond  float64 `json:"Bytes Written/Second"`
	CurrentTime         float64 `json:"Current Time"`
	CurrentRequests     int     `json:"Current Requests"`
	RequestsSecond      float64 `json:"Requests/Second"`
	ServerVersion       string  `json:"Server Version"`
	StartTime           float64 `json:"Start Time"`
	TotalBytesRead      int     `json:"Total Bytes Read"`
	TotalBytesWritten   int     `json:"Total Bytes Written"`
	TotalRequests       int     `json:"Total Requests"`
	TotalTime           float64 `json:"Total Time"`
	Uptime              float64 `json:"Uptime"`
	Requests            struct {
		BytesRead      int     `json:"Bytes Read"`
		BytesWritten   int     `json:"Bytes Written"`
		ResponseStatus string  `json:"Response Status"`
		StartTime      float64 `json:"Start Time"`
		EndTime        float64 `json:"End Time"`
		Client         string  `json:"Client"`
		ProcessingTime float64 `json:"Processing Time"`
		RequestLine    string  `json:"Request-Line"`
	} `json:"Requests"`
}

type HTTPServer struct {
	Enabled         bool    `json:"Enabled"`
	BindAddress     string  `json:"Bind Address"`
	RunTime         int     `json:"Run time"`
	Accepts         int     `json:"Accepts"`
	AcceptsSec      float64 `json:"Accepts/sec"`
	Queue           int     `json:"Queue"`
	Threads         int     `json:"Threads"`
	ThreadsIdle     int     `json:"Threads Idle"`
	SocketErrors    int     `json:"Socket Errors"`
	Requests        int     `json:"Requests"`
	BytesRead       int     `json:"Bytes Read"`
	BytesWritten    int     `json:"Bytes Written"`
	WorkTime        int     `json:"Work Time"`
	ReadThroughput  float64 `json:"Read Throughput"`
	WriteThroughput float64 `json:"Write Throughput"`
	WorkerThreads   map[string]struct {
		Requests        int     `json:"Requests"`
		BytesRead       int     `json:"Bytes Read"`
		BytesWritten    int     `json:"Bytes Written"`
		WorkTime        float64 `json:"Work Time"`
		ReadThroughput  float64 `json:"Read Throughput"`
		WriteThroughput float64 `json:"Write Throughput"`
	} `json:"Worker Threads"`
}

func (c *Client) Stats() (*StatsResponse, error) {
	data, err := c.doRequest("GET", "stats", nil)

	if err != nil {
		return nil, err
	}
	var tmp TMPStats
	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}

	var (
		app    Applications
		server HTTPServer
	)
	for key, val := range tmp {
		if key == CP_APPLICATIONS {
			td, _ := json.Marshal(val)

			if err := json.Unmarshal(td, &app); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(key, CP_HTTPSERVER_PREFIX) {
			td, _ := json.Marshal(val)

			if err := json.Unmarshal(td, &server); err != nil {
				return nil, err
			}
		}
	}

	stats := &StatsResponse{HTTPServer: server, Applications: app}

	return stats, nil
}
