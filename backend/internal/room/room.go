package room

import "errors"

var (
	ErrNotMaster   = errors.New("only the scrum master can do this")
	ErrNotCreator  = errors.New("only the room creator can do this")
	ErrNotFound    = errors.New("player not found")
	ErrInvalidCard = errors.New("invalid card value")
)

// CardValues is the ordered list of valid planning poker values.
var CardValues = []string{"0", "1", "2", "3", "5", "8", "13", "21", "40", "100", "?", "☕"}

var validCards map[string]bool

func init() {
	validCards = make(map[string]bool, len(CardValues))
	for _, v := range CardValues {
		validCards[v] = true
	}
}

type Player struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Vote     string `json:"vote,omitempty"`
	HasVoted bool   `json:"has_voted"`
}

type Room struct {
	ID        string
	CreatorID string
	MasterID  string
	Players   map[string]*Player
	Revealed  bool
	Round     int
}

func New(id, creatorID string) *Room {
	return &Room{
		ID:        id,
		CreatorID: creatorID,
		MasterID:  creatorID,
		Players:   make(map[string]*Player),
		Round:     1,
	}
}

func (r *Room) AddPlayer(id, name string) {
	r.Players[id] = &Player{ID: id, Name: name}
}

func (r *Room) RemovePlayer(id string) {
	delete(r.Players, id)
	if len(r.Players) == 0 {
		return
	}
	if r.CreatorID == id {
		for pid := range r.Players {
			r.CreatorID = pid
			break
		}
	}
	if r.MasterID == id {
		r.MasterID = r.CreatorID
	}
}

func (r *Room) Vote(playerID, value string) error {
	if !validCards[value] {
		return ErrInvalidCard
	}
	p, ok := r.Players[playerID]
	if !ok {
		return ErrNotFound
	}
	p.Vote = value
	p.HasVoted = true
	return nil
}

func (r *Room) Reveal(requesterID string) error {
	if r.MasterID != requesterID {
		return ErrNotMaster
	}
	r.Revealed = true
	return nil
}

func (r *Room) Reset(requesterID string) error {
	if r.MasterID != requesterID {
		return ErrNotMaster
	}
	r.Revealed = false
	r.Round++
	for _, p := range r.Players {
		p.Vote = ""
		p.HasVoted = false
	}
	return nil
}

func (r *Room) SetMaster(requesterID, targetID string) error {
	if r.CreatorID != requesterID {
		return ErrNotCreator
	}
	if _, ok := r.Players[targetID]; !ok {
		return ErrNotFound
	}
	r.MasterID = targetID
	return nil
}

// PlayerView is the client-facing representation of a player.
type PlayerView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Vote     string `json:"vote"`
	HasVoted bool   `json:"has_voted"`
}

// View is the client-facing representation of a room.
type View struct {
	ID        string       `json:"id"`
	CreatorID string       `json:"creator_id"`
	MasterID  string       `json:"master_id"`
	Players   []PlayerView `json:"players"`
	Revealed  bool         `json:"revealed"`
	Round     int          `json:"round"`
}

// ToView returns a serializable view, hiding individual votes until revealed.
func (r *Room) ToView() View {
	players := make([]PlayerView, 0, len(r.Players))
	for _, p := range r.Players {
		pv := PlayerView{
			ID:       p.ID,
			Name:     p.Name,
			HasVoted: p.HasVoted,
		}
		if r.Revealed {
			pv.Vote = p.Vote
		}
		players = append(players, pv)
	}
	return View{
		ID:        r.ID,
		CreatorID: r.CreatorID,
		MasterID:  r.MasterID,
		Players:   players,
		Revealed:  r.Revealed,
		Round:     r.Round,
	}
}
