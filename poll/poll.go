package poll

type Poll struct {
	ID        string `json:"_id"`
	Question  string
	CreatorID string
	Options   []string
}

type Vote struct {
	ID       string `json:"_id"`
	VoterID  string
	PollID   string
	VotedFor int
}

type Store interface {
	AddPoll(p Poll) error
	AddVote(v Vote) error
	GetResult(pollId string) (map[int]uint64, error)
	GetPoll(pollId string) (Poll, error)
	GetVote(voteId string) (Vote, error)
}

type StoreBackend interface {
	AddPoll(p Poll) error
	AddVote(v Vote) error
	GetPoll(pollId string) (Poll, error)
	GetVote(voteId string) (Vote, error)
	GetVotesForPoll(pollId string) ([]Vote, error)
	PollHasVoteFromVoter(pollID, voterID string) (bool, Vote, error)
}
