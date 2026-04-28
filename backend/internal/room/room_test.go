package room_test

import (
	"testing"

	"github.com/ericsantos/pokercards/internal/room"
)

func newRoom() *room.Room {
	r := room.New("room1", "p1")
	r.AddPlayer("p1", "Alice")
	return r
}

func TestAddPlayer(t *testing.T) {
	r := newRoom()
	r.AddPlayer("p2", "Bob")

	if len(r.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(r.Players))
	}
}

func TestRemovePlayer_TransfersCreator(t *testing.T) {
	r := newRoom()
	r.AddPlayer("p2", "Bob")

	r.RemovePlayer("p1")

	if r.CreatorID == "p1" {
		t.Error("creator should have been transferred after p1 left")
	}
	if r.CreatorID != "p2" {
		t.Errorf("expected creator p2, got %s", r.CreatorID)
	}
}

func TestRemovePlayer_TransfersMaster(t *testing.T) {
	r := newRoom()
	r.AddPlayer("p2", "Bob")
	if err := r.SetMaster("p1", "p2"); err != nil {
		t.Fatalf("SetMaster: %v", err)
	}

	r.RemovePlayer("p2")

	if r.MasterID == "p2" {
		t.Error("master should have been transferred after p2 left")
	}
}

func TestVote(t *testing.T) {
	r := newRoom()

	if err := r.Vote("p1", "8"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Players["p1"].HasVoted {
		t.Error("player should be marked as voted")
	}
}

func TestVote_InvalidCard(t *testing.T) {
	r := newRoom()
	if err := r.Vote("p1", "999"); err == nil {
		t.Error("expected error for invalid card")
	}
}

func TestVote_UnknownPlayer(t *testing.T) {
	r := newRoom()
	if err := r.Vote("unknown", "5"); err == nil {
		t.Error("expected error for unknown player")
	}
}

func TestReveal_RequiresMaster(t *testing.T) {
	r := newRoom()
	r.AddPlayer("p2", "Bob")

	if err := r.Reveal("p2"); err == nil {
		t.Error("non-master should not be able to reveal")
	}
	if err := r.Reveal("p1"); err != nil {
		t.Fatalf("master reveal failed: %v", err)
	}
	if !r.Revealed {
		t.Error("room should be revealed")
	}
}

func TestReset(t *testing.T) {
	r := newRoom()
	_ = r.Vote("p1", "5")
	_ = r.Reveal("p1")

	if err := r.Reset("p1"); err != nil {
		t.Fatalf("reset failed: %v", err)
	}
	if r.Revealed {
		t.Error("room should not be revealed after reset")
	}
	if r.Players["p1"].HasVoted {
		t.Error("votes should be cleared after reset")
	}
	if r.Round != 2 {
		t.Errorf("expected round 2, got %d", r.Round)
	}
}

func TestSetMaster_RequiresCreator(t *testing.T) {
	r := newRoom()
	r.AddPlayer("p2", "Bob")

	if err := r.SetMaster("p2", "p2"); err == nil {
		t.Error("non-creator should not be able to set master")
	}
	if err := r.SetMaster("p1", "p2"); err != nil {
		t.Fatalf("set master failed: %v", err)
	}
	if r.MasterID != "p2" {
		t.Errorf("master should be p2, got %s", r.MasterID)
	}
}

func TestSetMaster_UnknownTarget(t *testing.T) {
	r := newRoom()
	if err := r.SetMaster("p1", "unknown"); err == nil {
		t.Error("expected error for unknown target player")
	}
}

func TestToView_HidesVotesBeforeReveal(t *testing.T) {
	r := newRoom()
	_ = r.Vote("p1", "8")

	v := r.ToView()
	for _, p := range v.Players {
		if p.Vote != "" {
			t.Error("vote should be hidden before reveal")
		}
		if !p.HasVoted {
			t.Error("has_voted should be true")
		}
	}
}

func TestToView_ShowsVotesAfterReveal(t *testing.T) {
	r := newRoom()
	_ = r.Vote("p1", "8")
	_ = r.Reveal("p1")

	v := r.ToView()
	for _, p := range v.Players {
		if p.Vote != "8" {
			t.Errorf("expected vote '8', got %q", p.Vote)
		}
	}
}
