package salt

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Event represents a SaltStack event from the event bus
type Event struct {
	Tag       string    `json:"tag"`
	Data      any       `json:"data"`
	Stamp     time.Time `json:"stamp"`
	AlterTime time.Time `json:"alter_time"`
}

// EventOption configures event stream parameters
type EventOption func(*eventConfig)

type eventConfig struct {
	format string
}

// WithEventFormat sets the response format for events (default: "application/json")
// Accepts: "application/json" or "application/x-yaml"
func WithEventFormat(format string) EventOption {
	return func(ec *eventConfig) {
		ec.format = format
	}
}

// Events streams SaltStack events from the event bus using Server-Sent Events (SSE)
// https://docs.saltproject.io/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#events
//
// The Events endpoint is designed for long-running HTTP connections and exposes
// Salt's event bus, which includes job returns and other system events.
//
// Events are streamed in real-time as they occur. The method blocks until
// the context is cancelled.
//
// Example:
//
//	client := salt.NewClient(...)
//	if err := client.Login(ctx); err != nil {
//		return err
//	}
//
//	eventChan, err := client.Events(ctx)
//	if err != nil {
//		return err
//	}
//
//	for event := range eventChan {
//		fmt.Printf("Event: %+v\n", event)
//	}
func (c *Client) Events(ctx context.Context, opts ...EventOption) (<-chan Event, error) {
	config := &eventConfig{
		format: "application/json",
	}
	for _, opt := range opts {
		opt(config)
	}

	// Create a channel for streaming events
	eventChan := make(chan Event, 100)

	// Create request with authentication
	url := c.endpoint + "/events"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		close(eventChan)
		return eventChan, err
	}

	req.Header.Set("Accept", config.format)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if c.token != "" {
		req.Header.Set("X-Auth-Token", c.token)
	}

	// Start streaming in a goroutine
	go func() {
		defer close(eventChan)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		// Parse SSE (Server-Sent Events) format
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// SSE format: "data: {json}"
			if jsonData, ok := strings.CutPrefix(line, "data:"); ok {
				jsonData = strings.TrimSpace(jsonData)

				// Parse event from JSON
				var event Event
				if err := parseEventFromJSON(jsonData, &event); err != nil {
					// Try parsing as gjson for more flexibility
					if parsedEvent := parseEventFromGJSON(jsonData); parsedEvent != nil {
						eventChan <- *parsedEvent
					}
					continue
				}

				eventChan <- event
			}
		}

		if err := scanner.Err(); err != nil {
			// Context cancelled is expected, not an error
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return eventChan, nil
}

// parseEventFromJSON parses an event from JSON string
func parseEventFromJSON(jsonData string, event *Event) error {
	result := gjson.Parse(jsonData)

	// Parse tag
	if tag := result.Get("tag"); tag.Exists() {
		event.Tag = tag.String()
	}

	// Parse data as raw JSON
	if data := result.Get("data"); data.Exists() {
		event.Data = data.Value()
	}

	// Parse timestamp
	if stamp := result.Get("stamp"); stamp.Exists() {
		// Try parsing timestamp (various formats)
		if stamp.Type == gjson.String {
			// Parse as string timestamp
			t, err := time.Parse("2006-01-02T15:04:05.000000", stamp.String())
			if err == nil {
				event.Stamp = t
			}
		} else if stamp.Type == gjson.Number {
			// Parse as Unix timestamp
			event.Stamp = time.Unix(stamp.Int(), 0)
		}
	}

	// Parse alter_time
	if alterTime := result.Get("alter_time"); alterTime.Exists() {
		t, err := time.Parse("2006-01-02T15:04:05.000000", alterTime.String())
		if err == nil {
			event.AlterTime = t
		}
	}

	return nil
}

// parseEventFromGJSON parses event using gjson for more flexible parsing
func parseEventFromGJSON(jsonData string) *Event {
	result := gjson.Parse(jsonData)

	if !result.Exists() {
		return nil
	}

	event := &Event{
		Tag:  result.Get("tag").String(),
		Data: result.Get("data").Value(),
	}

	// Parse timestamp
	if stamp := result.Get("stamp"); stamp.Exists() {
		if stamp.Type == gjson.Number {
			event.Stamp = time.Unix(stamp.Int(), 0)
		} else {
			t, err := time.Parse("2006-01-02T15:04:05.000000", stamp.String())
			if err == nil {
				event.Stamp = t
			}
		}
	}

	return event
}

// EventData extracts typed data from an event
func (e *Event) EventData() map[string]any {
	data, ok := e.Data.(map[string]any)
	if !ok {
		return nil
	}
	return data
}

// String returns a string representation of the event
func (e *Event) String() string {
	return fmt.Sprintf("Event{Tag: %s, Stamp: %s, Data: %v}", e.Tag, e.Stamp, e.Data)
}

// EventsReader provides an alternative interface for reading events with an io.Reader style
type EventsReader struct {
	reader *bufio.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

// NewEventsReader creates a new EventsReader for custom event handling
func (c *Client) NewEventsReader(ctx context.Context, opts ...EventOption) (*EventsReader, error) {
	config := &eventConfig{
		format: "application/json",
	}
	for _, opt := range opts {
		opt(config)
	}

	url := c.endpoint + "/events"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", config.format)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if c.token != "" {
		req.Header.Set("X-Auth-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("events endpoint returned status %d", resp.StatusCode)
	}

	return &EventsReader{
		reader: bufio.NewReader(resp.Body),
		ctx:    ctx,
	}, nil
}

// ReadEvent reads the next event from the stream
// Returns io.EOF when the stream is closed or context is cancelled
func (er *EventsReader) ReadEvent() (*Event, error) {
	for {
		select {
		case <-er.ctx.Done():
			return nil, er.ctx.Err()
		default:
			line, err := er.reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil, io.EOF
				}
				return nil, err
			}

			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			if jsonData, ok := strings.CutPrefix(line, "data:"); ok {
				jsonData = strings.TrimSpace(jsonData)
				event := parseEventFromGJSON(jsonData)
				if event != nil {
					return event, nil
				}
			}
		}
	}
}

// Close closes the event stream
func (er *EventsReader) Close() error {
	if er.cancel != nil {
		er.cancel()
	}
	return nil
}
