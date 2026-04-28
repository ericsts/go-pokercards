package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ericsantos/pokercards/internal/room"
	"github.com/ericsantos/pokercards/internal/session"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Allow all origins for development. In production the nginx proxy
	// handles origin validation at the edge.
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	store *session.Store
}

func New(store *session.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/cards", h.listCards)
	mux.HandleFunc("POST /api/rooms", h.createRoom)
	mux.HandleFunc("GET /api/rooms/{id}", h.getRoom)
	mux.HandleFunc("GET /api/rooms/{id}/ws", h.handleWS)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listCards(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"cards": room.CardValues,
	})
}

func (h *Handler) createRoom(w http.ResponseWriter, _ *http.Request) {
	roomID := uuid.New().String()[:8]
	creatorID := uuid.New().String()
	h.store.Create(roomID, creatorID)
	writeJSON(w, http.StatusCreated, map[string]string{
		"room_id":    roomID,
		"creator_id": creatorID,
	})
}

func (h *Handler) getRoom(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sess, ok := h.store.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "room not found"})
		return
	}
	sess.RLock()
	view := sess.Room.ToView()
	sess.RUnlock()
	writeJSON(w, http.StatusOK, view)
}

func (h *Handler) handleWS(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name query param is required", http.StatusBadRequest)
		return
	}

	sess, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "err", err)
		return
	}

	// Use a provided player_id for reconnection (e.g. the room creator).
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		playerID = uuid.New().String()
	}

	// Add player to the room data (locked).
	sess.Lock()
	sess.Room.AddPlayer(playerID, name)
	initView := sess.Room.ToView()
	sess.Unlock()

	// Prepare the client (starts its write goroutine) before registering
	// so the init message is guaranteed to arrive before any broadcast.
	client := sess.PrepareClient(conn)
	client.Send(map[string]interface{}{
		"type":      "init",
		"player_id": playerID,
		"room":      initView,
	})

	sess.RegisterClient(playerID, client)
	sess.BroadcastState()

	defer func() {
		sess.Lock()
		sess.Room.RemovePlayer(playerID)
		isEmpty := len(sess.Room.Players) == 0
		sess.Unlock()

		sess.RemoveClient(playerID)

		if isEmpty {
			h.store.Delete(id)
		} else {
			sess.BroadcastState()
		}
	}()

	for {
		var msg session.ClientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		h.handleMessage(sess, playerID, msg)
	}
}

func (h *Handler) handleMessage(sess *session.Session, playerID string, msg session.ClientMessage) {
	var err error

	sess.Lock()
	switch msg.Action {
	case "vote":
		err = sess.Room.Vote(playerID, msg.Value)
	case "reveal":
		err = sess.Room.Reveal(playerID)
	case "reset":
		err = sess.Room.Reset(playerID)
	case "set_master":
		err = sess.Room.SetMaster(playerID, msg.PlayerID)
	default:
		sess.Unlock()
		return
	}
	sess.Unlock()

	if err != nil {
		sess.SendError(playerID, err.Error())
		return
	}
	sess.BroadcastState()
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
