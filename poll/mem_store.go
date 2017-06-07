package poll

import (
	"sync"
)

type InMemoryStore struct {
	pollStore map[string]Poll
	voteStore map[string][]Vote
	lock      sync.Mutex
}

func NewInMemoryStore() StoreBackend {
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

func (s *InMemoryStore) GetVotesForPoll(pollId string) ([]Vote, error) {
	s.lock.Lock()
	votes := s.voteStore[pollId]
	s.lock.Unlock()
	return votes, nil
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
