package poll

import (
	"sync"
)

type Poll struct {
	ID        string `json:"_id"`
	Question  string
	CreatorID string
	Options   []string
}

func NewPoll(id, question, creatorID string, options []string) Poll {
	return Poll{id, question, creatorID, options}
}

type Vote struct {
	ID       string `json:"_id"`
	VoterID  string
	PollID   string
	VotedFor string
}

func NewVote(id, voterID, pollID, votedFor string) Vote {
	return Vote{id, voterID, pollID, votedFor}
}

type Store interface {
	AddPoll(p Poll) error
	AddVote(v Vote) error
	GetResult(pollId string) (map[string]uint64, error)
	GetPoll(pollId string) (Poll, error)
	GetVote(voteId string) (Vote, error)
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

func (s *InMemoryStore) AddPoll(p Poll) error {
	s.lock.Lock()
	s.pollStore[p.ID] = p
	s.voteStore[p.ID] = make([]Vote, 0, 20)
	s.lock.Unlock()
	return nil
}

func (s *InMemoryStore) AddVote(v Vote) error {
	s.lock.Lock()
	oldVotes := s.voteStore[v.PollID]
	s.voteStore[v.PollID] = append(oldVotes, v)
	s.lock.Unlock()
	return nil
}

func (s *InMemoryStore) GetResult(pollId string) (map[string]uint64, error) {
	result := make(map[string]uint64)
	s.lock.Lock()
	votes := s.voteStore[pollId]
	for _, vote := range votes {
		result[vote.VotedFor]++
	}
	s.lock.Unlock()
	return result, nil
}

func (s *InMemoryStore) GetPoll(pollId string) (Poll, error) {
	s.lock.Lock()
	poll := s.pollStore[pollId]
	s.lock.Unlock()
	return poll, nil
}

func (s *InMemoryStore) GetVote(voteId string) (Vote, error) {
	var foundVote Vote
	s.lock.Lock()
	for _, votes := range s.voteStore {
		for _, vote := range votes {
			if vote.ID == voteId {
				foundVote = vote
				break
			}
		}
	}
	s.lock.Unlock()
	return foundVote, nil
}
