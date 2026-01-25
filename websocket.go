package salt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

// WSOption configures WebSocket connection parameters
type WSOption func(*wsConfig)

type wsConfig struct {
	formatEvents bool
	headers      http.Header
}

// WithWSFormatEvents enables formatted event output instead of raw JSON
// When enabled, events are formatted with minion/job information
func WithWSFormatEvents(format bool) WSOption {
	return func(wc *wsConfig) {
		wc.formatEvents = format
	}
}

// WithWSHeader sets custom headers for WebSocket connection
func WithWSHeader(headers http.Header) WSOption {
	return func(wc *wsConfig) {
		wc.headers = headers
	}
}

// WebSocketClient represents a WebSocket connection to SaltStack
type WebSocketClient struct {
	conn         *websocket.Conn
	endpoint     string
	token        string
	formatEvents bool
	headers      http.Header
	eventChan    chan []byte
	done         chan struct{}
	mu           sync.RWMutex
	closed       bool
}

// WebSocket represents a WebSocket message
type WebSocketMessage struct {
	Type int    `json:"type"`
	Data []byte `json:"data"`
}

// Connect establishes a WebSocket connection to SaltStack event bus
// https://docs.saltproject.io/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#ws
//
// The WebSocket endpoint provides real-time bidirectional access to Salt's
// event bus. Unlike the SSE Events endpoint, WebSocket allows
// sending messages (e.g., "ready" signal) and receiving events.
//
// Connection format:
// - ws://host:port/ws for non-SSL
// - wss://host:port/ws for SSL
// - Token can be in URL: /ws/<token> or in X-Auth-Token header
//
// Example:
//
//	client := salt.NewClient(...)
//	if err := client.Login(ctx); err != nil {
//		return err
//	}
//
//	wsClient, err := client.NewWebSocket(ctx)
//	if err != nil {
//		return err
//	}
//	defer wsClient.Close()
//
//	// Send ready signal to start receiving events
//	if err := wsClient.SendReady(ctx); err != nil {
//		return err
//	}
//
//	// Receive events
//	for {
//		msg, err := wsClient.Receive(ctx)
//		if err != nil {
//			break
//		}
//		fmt.Printf("Event: %s\n", string(msg))
//	}
func (c *Client) NewWebSocket(ctx context.Context, opts ...WSOption) (*WebSocketClient, error) {
	config := &wsConfig{
		formatEvents: false,
		headers:      make(http.Header),
	}
	for _, opt := range opts {
		opt(config)
	}

	// Build WebSocket URL
	wsURL := fmt.Sprintf("%s/ws", c.endpoint)
	if c.token != "" {
		wsURL = fmt.Sprintf("%s/ws/%s", c.endpoint, c.token)
	}

	// Convert http:// to ws://, https:// to wss://
	wsURL = strings.Replace(wsURL, "http://", "ws://", 1)
	wsURL = strings.Replace(wsURL, "https://", "wss://", 1)

	// Set authentication header
	if c.token != "" && len(config.headers.Get("X-Auth-Token")) == 0 {
		config.headers.Set("X-Auth-Token", c.token)
	}

	// Create request with custom headers
	reqHeader := http.Header{}
	for k, v := range config.headers {
		reqHeader[k] = v
	}

	// Parse URL
	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, parsedURL.String(), reqHeader)
	if err != nil {
		return nil, fmt.Errorf("WebSocket connection failed: %w", err)
	}

	return &WebSocketClient{
		conn:         conn,
		endpoint:     c.endpoint,
		token:        c.token,
		formatEvents: config.formatEvents,
		headers:      config.headers,
		eventChan:    make(chan []byte, 100),
		done:         make(chan struct{}),
	}, nil
}

// SendReady sends the required "websocket client ready" signal
// This signal must be sent before the server starts sending events
func (wsc *WebSocketClient) SendReady(ctx context.Context) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.closed {
		return fmt.Errorf("WebSocket connection is closed")
	}

	return wsc.conn.WriteMessage(websocket.TextMessage, []byte("websocket client ready"))
}

// Send sends a message to the WebSocket server
func (wsc *WebSocketClient) Send(ctx context.Context, messageType int, data []byte) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.closed {
		return fmt.Errorf("WebSocket connection is closed")
	}

	return wsc.conn.WriteMessage(messageType, data)
}

// Receive reads the next message from the WebSocket server
// Blocks until a message is received, context is cancelled, or connection is closed
func (wsc *WebSocketClient) Receive(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-wsc.done:
		return nil, fmt.Errorf("WebSocket connection is closed")
	case msg, ok := <-wsc.eventChan:
		if !ok {
			return nil, fmt.Errorf("event channel closed")
		}
		return msg, nil
	}
}

// ReceiveEvent reads and parses the next event from the WebSocket
// Returns a parsed Event struct or an error
func (wsc *WebSocketClient) ReceiveEvent(ctx context.Context) (*Event, error) {
	data, err := wsc.Receive(ctx)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(data)
	if !result.Exists() {
		return nil, fmt.Errorf("invalid event data")
	}

	event := &Event{
		Tag:  result.Get("tag").String(),
		Data: result.Get("data").Value(),
	}

	// Parse timestamp
	if stamp := result.Get("stamp"); stamp.Exists() {
		switch stamp.Type {
		case gjson.Number:
			event.Stamp = time.Unix(stamp.Int(), 0)
		default:
			t, err := time.Parse("2006-01-02T15:04:05.000000", stamp.String())
			if err == nil {
				event.Stamp = t
			}
		}
	}

	// Parse alter_time
	if alterTime := result.Get("alter_time"); alterTime.Exists() {
		t, err := time.Parse("2006-01-02T15:04:05.000000", alterTime.String())
		if err == nil {
			event.AlterTime = t
		}
	}

	return event, nil
}

// Close closes the WebSocket connection gracefully
func (wsc *WebSocketClient) Close() error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.closed {
		return nil
	}

	wsc.closed = true
	close(wsc.eventChan)
	close(wsc.done)

	return wsc.conn.Close()
}

// IsClosed returns true if the WebSocket connection is closed
func (wsc *WebSocketClient) IsClosed() bool {
	wsc.mu.RLock()
	defer wsc.mu.RUnlock()
	return wsc.closed
}

// ReceiveJSON reads the next message and unmarshals it into the provided struct
func (wsc *WebSocketClient) ReceiveJSON(ctx context.Context, v any) error {
	data, err := wsc.Receive(ctx)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// StartListener starts a goroutine to receive events automatically
// Events are sent to the provided channel
func (wsc *WebSocketClient) StartListener(ctx context.Context) (<-chan []byte, error) {
	wsc.mu.Lock()
	if wsc.closed {
		wsc.mu.Unlock()
		return nil, fmt.Errorf("WebSocket connection is closed")
	}
	wsc.mu.Unlock()

	outChan := make(chan []byte, 100)

	go func() {
		defer close(outChan)
		defer wsc.Close()

		for {
			messageType, data, err := wsc.conn.ReadMessage()
			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {
				select {
				case outChan <- data:
				case <-ctx.Done():
					return
				case <-wsc.done:
					return
				}
			}
		}
	}()

	return outChan, nil
}

// StartEventListener starts receiving events automatically and returns them as Event structs
func (wsc *WebSocketClient) StartEventListener(ctx context.Context) (<-chan Event, error) {
	rawChan, err := wsc.StartListener(ctx)
	if err != nil {
		return nil, err
	}

	eventChan := make(chan Event, 100)

	go func() {
		defer close(eventChan)

		for raw := range rawChan {
			result := gjson.ParseBytes(raw)
			if !result.Exists() {
				continue
			}

			event := &Event{
				Tag:  result.Get("tag").String(),
				Data: result.Get("data").Value(),
			}

			// Parse timestamp
			if stamp := result.Get("stamp"); stamp.Exists() {
				switch stamp.Type {
				case gjson.Number:
					event.Stamp = time.Unix(stamp.Int(), 0)
				default:
					t, err := time.Parse("2006-01-02T15:04:05.000000", stamp.String())
					if err == nil {
						event.Stamp = t
					}
				}
			}

			// Parse alter_time
			if alterTime := result.Get("alter_time"); alterTime.Exists() {
				t, err := time.Parse("2006-01-02T15:04:05.000000", alterTime.String())
				if err == nil {
					event.AlterTime = t
				}
			}

			select {
			case eventChan <- *event:
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventChan, nil
}
