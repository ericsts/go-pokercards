package session_test

import (
	"testing"

	"github.com/ericsantos/pokercards/internal/session"
)

func TestStore_CreateAndGet(t *testing.T) {
	s := session.NewStore()
	sess := s.Create("r1", "p1")

	if sess == nil {
		t.Fatal("Create should return a non-nil session")
	}
	if sess.Room.ID != "r1" {
		t.Errorf("expected room ID 'r1', got %q", sess.Room.ID)
	}
	if sess.Room.CreatorID != "p1" {
		t.Errorf("expected creator 'p1', got %q", sess.Room.CreatorID)
	}

	got, ok := s.Get("r1")
	if !ok {
		t.Fatal("Get should find the session after Create")
	}
	if got != sess {
		t.Error("Get should return the same session pointer")
	}
}

func TestStore_GetMissing(t *testing.T) {
	s := session.NewStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("Get should return false for nonexistent room")
	}
}

func TestStore_Delete(t *testing.T) {
	s := session.NewStore()
	s.Create("r1", "p1")
	s.Delete("r1")

	if _, ok := s.Get("r1"); ok {
		t.Error("Get should return false after Delete")
	}
	if s.Count() != 0 {
		t.Errorf("Count should be 0 after Delete, got %d", s.Count())
	}
}

func TestStore_Count(t *testing.T) {
	s := session.NewStore()
	if s.Count() != 0 {
		t.Error("initial count should be 0")
	}

	s.Create("r1", "p1")
	s.Create("r2", "p2")
	s.Create("r3", "p3")

	if s.Count() != 3 {
		t.Errorf("expected count 3, got %d", s.Count())
	}

	s.Delete("r2")
	if s.Count() != 2 {
		t.Errorf("expected count 2 after delete, got %d", s.Count())
	}

	if _, ok := s.Get("r1"); !ok {
		t.Error("r1 should still exist after deleting r2")
	}
	if _, ok := s.Get("r3"); !ok {
		t.Error("r3 should still exist after deleting r2")
	}
}
