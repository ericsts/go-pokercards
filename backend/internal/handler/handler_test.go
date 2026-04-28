package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ericsantos/pokercards/internal/handler"
	"github.com/ericsantos/pokercards/internal/session"
)

func newTestServer() *httptest.Server {
	store := session.NewStore()
	h := handler.New(store)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func TestHealth(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestListCards(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/cards")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	cards, ok := body["cards"]
	if !ok || cards == nil {
		t.Error("response should have 'cards' key")
	}
}

func TestCreateRoom(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/rooms", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["room_id"] == "" {
		t.Error("response should include room_id")
	}
	if body["creator_id"] == "" {
		t.Error("response should include creator_id")
	}
}

func TestGetRoom_NotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/rooms/doesnotexist")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetRoom_Found(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// Create a room first.
	resp, err := http.Post(srv.URL+"/api/rooms", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var created map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	// Now fetch it.
	resp2, err := http.Get(srv.URL + "/api/rooms/" + created["room_id"])
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	var room map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&room); err != nil {
		t.Fatal(err)
	}
	if room["id"] != created["room_id"] {
		t.Errorf("room ID mismatch: %v", room["id"])
	}
}
