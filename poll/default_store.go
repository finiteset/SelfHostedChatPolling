package poll

type DefaultStore struct {
	backend StoreBackend
}

func NewDefaultStore(backend StoreBackend) Store {
	return &DefaultStore{backend}
}

func (s *DefaultStore) AddPoll(p Poll) error {
	return s.backend.AddPoll(p)
}

func (s *DefaultStore) AddVote(v Vote) error {
	return s.backend.AddVote(v)
}

func (s *DefaultStore) GetResult(pollId string) (map[string]uint64, error) {
	result := make(map[string]uint64)
	votes, err := s.backend.GetVotesForPoll(pollId)
	if err != nil {
		return nil, err
	}
	for _, vote := range votes {
		result[vote.VotedFor]++
	}
	return result, nil
}

func (s *DefaultStore) GetPoll(pollId string) (Poll, error) {
	return s.backend.GetPoll(pollId)
}

func (s *DefaultStore) GetVote(voteId string) (Vote, error) {
	return s.backend.GetVote(voteId)
}
