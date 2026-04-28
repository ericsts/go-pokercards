package session

import "github.com/gorilla/websocket"

// Client wraps a WebSocket connection with a dedicated writer goroutine,
// making it safe to send messages from multiple goroutines.
type Client struct {
	send chan interface{}
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *Client {
	c := &Client{
		send: make(chan interface{}, 32),
		conn: conn,
	}
	go c.writePump()
	return c
}

// writePump drains the send channel and writes each message to the connection.
// Gorilla WebSocket requires all writes from a single goroutine.
func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			return
		}
	}
}

// Send queues a message for delivery. Drops the message if the buffer is full.
func (c *Client) Send(msg interface{}) {
	select {
	case c.send <- msg:
	default:
	}
}

// Close drains and closes the send channel, stopping the writer goroutine.
func (c *Client) Close() {
	close(c.send)
}
