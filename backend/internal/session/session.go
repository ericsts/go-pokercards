package session

import (
	"sync"

	"github.com/ericsantos/pokercards/internal/room"
	"github.com/gorilla/websocket"
)

// ServerMessage is the envelope the server sends to clients.
type ServerMessage struct {
	Type    string     `json:"type"`
	Room    *room.View `json:"room,omitempty"`
	Message string     `json:"message,omitempty"`
}

// ClientMessage is the envelope received from clients.
type ClientMessage struct {
	Action   string `json:"action"`
	Value    string `json:"value,omitempty"`
	PlayerID string `json:"player_id,omitempty"`
}

// Session groups a Room with its active WebSocket clients.
// The embedded RWMutex protects both Room fields and the clients map.
type Session struct {
	sync.RWMutex
	Room    *room.Room
	clients map[string]*Client
}

func newSession(r *room.Room) *Session {
	return &Session{
		Room:    r,
		clients: make(map[string]*Client),
	}
}

// PrepareClient creates a Client (starts its write pump) without registering it.
// Call RegisterClient after sending the init message so the client receives
// its own init before any broadcast.
func (s *Session) PrepareClient(conn *websocket.Conn) *Client {
	return newClient(conn)
}

// RegisterClient adds a prepared client to the session so it receives broadcasts.
func (s *Session) RegisterClient(playerID string, c *Client) {
	s.Lock()
	s.clients[playerID] = c
	s.Unlock()
}

// RemoveClient closes and unregisters a client.
func (s *Session) RemoveClient(playerID string) {
	s.Lock()
	if c, ok := s.clients[playerID]; ok {
		c.Close()
		delete(s.clients, playerID)
	}
	s.Unlock()
}

// ClientCount returns the number of connected clients.
func (s *Session) ClientCount() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.clients)
}

// BroadcastState sends the current room state to all connected clients.
func (s *Session) BroadcastState() {
	s.RLock()
	view := s.Room.ToView()
	clients := make([]*Client, 0, len(s.clients))
	for _, c := range s.clients {
		clients = append(clients, c)
	}
	s.RUnlock()

	msg := ServerMessage{Type: "state", Room: &view}
	for _, c := range clients {
		c.Send(msg)
	}
}

// SendError sends an error message to a single player.
func (s *Session) SendError(playerID, errMsg string) {
	s.RLock()
	c, ok := s.clients[playerID]
	s.RUnlock()
	if ok {
		c.Send(ServerMessage{Type: "error", Message: errMsg})
	}
}
