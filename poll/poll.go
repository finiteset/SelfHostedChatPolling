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
	PollHasVoteFromVoter(pollID, voterID string) (bool, error)
}
