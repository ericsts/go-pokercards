package session

import (
	"sync"

	"github.com/ericsantos/pokercards/internal/room"
)

// Store manages all active game sessions.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]*Session)}
}

// Create creates a new session for the given room and creator IDs.
func (s *Store) Create(roomID, creatorID string) *Session {
	r := room.New(roomID, creatorID)
	sess := newSession(r)
	s.mu.Lock()
	s.sessions[roomID] = sess
	s.mu.Unlock()
	return sess
}

// Get retrieves a session by room ID.
func (s *Store) Get(roomID string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[roomID]
	return sess, ok
}

// Delete removes a session.
func (s *Store) Delete(roomID string) {
	s.mu.Lock()
	delete(s.sessions, roomID)
	s.mu.Unlock()
}

// Count returns the number of active sessions.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}
