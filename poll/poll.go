package poll

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

type StoreBackend interface {
	AddPoll(p Poll) error
	AddVote(v Vote) error
	GetPoll(pollId string) (Poll, error)
	GetVote(voteId string) (Vote, error)
	GetVotesForPoll(pollId string) ([]Vote, error)
}

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
