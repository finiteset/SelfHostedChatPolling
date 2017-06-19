package memstore

import (
	"markusreschke.name/selfhostedchatpolling/poll"
	"sync"
)

type InMemoryStore struct {
	pollStore map[string]poll.Poll
	voteStore map[string][]poll.Vote
	lock      sync.Mutex
}

func NewInMemoryStoreBackend() poll.StoreBackend {
	store := new(InMemoryStore)
	store.pollStore = make(map[string]poll.Poll)
	store.voteStore = make(map[string][]poll.Vote)
	return store
}

func (s *InMemoryStore) PollHasVoteFromVoter(pollID, voterID string) (bool, poll.Vote, error) {
	s.lock.Lock()
	voteFound := false
	foundVote := poll.Vote{}
	for _, vote := range s.voteStore[pollID] {
		if vote.VoterID == voterID {
			voteFound = true
			foundVote = vote
			break
		}
	}
	s.lock.Unlock()
	return voteFound, foundVote, nil
}

func (s *InMemoryStore) AddPoll(p poll.Poll) error {
	s.lock.Lock()
	s.pollStore[p.ID] = p
	s.voteStore[p.ID] = make([]poll.Vote, 0, 20)
	s.lock.Unlock()
	return nil
}

func (s *InMemoryStore) AddVote(v poll.Vote) error {
	s.lock.Lock()
	oldVotes := s.voteStore[v.PollID]
	s.voteStore[v.PollID] = append(oldVotes, v)
	s.lock.Unlock()
	return nil
}

func (s *InMemoryStore) GetVotesForPoll(pollId string) ([]poll.Vote, error) {
	s.lock.Lock()
	votes := s.voteStore[pollId]
	s.lock.Unlock()
	return votes, nil
}

func (s *InMemoryStore) GetPoll(pollId string) (poll.Poll, error) {
	s.lock.Lock()
	poll := s.pollStore[pollId]
	s.lock.Unlock()
	return poll, nil
}

func (s *InMemoryStore) GetVote(voteId string) (poll.Vote, error) {
	var foundVote poll.Vote
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
