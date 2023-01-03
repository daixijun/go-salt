package salt

import (
	"context"
	"encoding/json"
	"strings"
)

type TMPStats map[string]interface{}

const (
	CP_APPLICATIONS      = "CherryPy Applications"
	CP_HTTPSERVER_PREFIX = "Cheroot HTTPServer"
)

type Stats struct {
	Applications Applications `json:"CherryPy Applications"`
	HTTPServer   HTTPServer   `json:"Cheroot HTTPServer"`
}

func (c *Client) Stats(ctx context.Context) (*Stats, error) {
	data, err := c.get(ctx, "Stats")

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

	ret := &Stats{HTTPServer: server, Applications: app}

	return ret, nil
}
