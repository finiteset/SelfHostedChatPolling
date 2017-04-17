package poll

import "sync"

type Poll interface {
	Question() string
	ID() string
	CreatorID() string
}

type Vote interface {
	ID() string
	VoterID() string
	PollID() string
	VotedFor() string
}

type SimplePoll struct {
	id        string
	question  string
	creatorID string
}

func NewPoll(id, question, creatorID string) Poll { return SimplePoll{id, question, creatorID} }
func (p SimplePoll) ID() string                   { return p.id }
func (p SimplePoll) Question() string             { return p.question }
func (p SimplePoll) CreatorID() string            { return p.creatorID }

type SimpleVote struct {
	id       string
	voterID  string
	pollID   string
	votedFor string
}

func NewSimpleVote(id, voterID, pollID, votedFor string) Vote {
	return SimpleVote{id, voterID, pollID, votedFor}
}
func (v SimpleVote) ID() string       { return v.id }
func (v SimpleVote) VoterID() string  { return v.voterID }
func (v SimpleVote) PollID() string   { return v.pollID }
func (v SimpleVote) VotedFor() string { return v.votedFor }

type Store interface {
	AddPoll(p Poll)
	AddVote(v Vote)
	GetResult(p Poll) map[string]uint64
}

type InMemoryStore struct {
	pollStore map[string]Poll
	voteStore map[string][]Vote
	lock      sync.Mutex
}

func NewInMemoryStore() Store {
	store := new(InMemoryStore)
	store.pollStore = make(map[string]Poll)
	store.voteStore = make(map[string][]Vote)
	return store
}

func (s *InMemoryStore) AddPoll(p Poll) {
	s.lock.Lock()
	s.pollStore[p.ID()] = p
	s.voteStore[p.ID()] = make([]Vote, 20)
	s.lock.Unlock()
}
func (s *InMemoryStore) AddVote(v Vote) {
	s.lock.Lock()
	oldVotes := s.voteStore[v.PollID()]
	s.voteStore[v.PollID()] = append(oldVotes, v)
	s.lock.Unlock()
}
func (s *InMemoryStore) GetResult(p Poll) map[string]uint64 {
	result := make(map[string]uint64)
	s.lock.Lock()
	votes := s.voteStore[p.ID()]
	for _, vote := range votes {
		result[vote.VotedFor()]++
	}
	s.lock.Unlock()
	return result
}
