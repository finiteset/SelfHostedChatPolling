package poll

import (
	"fmt"

	"github.com/pkg/errors"
)

type DefaultStore struct {
	backend StoreBackend
}

var (
	ErrNoDetailsForAnonymousPoll = errors.New("Can't fetch vote details for anonymous polls!")
	ErrInvalidChoice             = errors.New("Invalid option choice!")
)

func NewDefaultStore(backend StoreBackend) Store {
	return &DefaultStore{backend}
}

func (s *DefaultStore) AddPoll(p Poll) error {
	return s.backend.AddPoll(p)
}

func (s *DefaultStore) AddVote(v Vote) error {
	isValidChoice, err := s.votedForValidOption(v)
	if err != nil {
		return err
	}
	if !isValidChoice {
		return errors.Wrap(ErrInvalidChoice, fmt.Sprintf("Voter %s voted for invalid choice %s", v.VoterID, v.VotedFor))
	}
	hasVotedAlready, previousVote, err := s.backend.PollHasVoteFromVoter(v.PollID, v.VoterID)
	if err != nil {
		return err
	}
	if hasVotedAlready {
		err := s.backend.RemoveVote(previousVote.ID)
		if err != nil {
			return err
		}
	}
	return s.backend.AddVote(v)
}

func (s *DefaultStore) votedForValidOption(v Vote) (bool, error) {
	pollForVote, err := s.backend.GetPoll(v.PollID)
	if err != nil {
		return false, err
	}
	isValidChoice := v.VotedFor >= 0 && v.VotedFor < len(pollForVote.Options)
	return isValidChoice, nil
}

func (s *DefaultStore) GetResult(pollId string) (map[int]uint64, error) {
	result := make(map[int]uint64)
	votes, err := s.backend.GetVotesForPoll(pollId)
	if err != nil {
		return nil, err
	}
	for _, vote := range votes {
		result[vote.VotedFor]++
	}
	return result, nil
}

func (s *DefaultStore) GetVoteDetails(pollId string) (map[string][]string, error) {
	result := make(map[string][]string)
	pollForId, err := s.backend.GetPoll(pollId)
	if err != nil {
		return nil, err
	}
	if pollForId.Anonymous {
		return nil, ErrNoDetailsForAnonymousPoll
	}
	votes, err := s.backend.GetVotesForPoll(pollId)
	if err != nil {
		return nil, err
	}
	for _, optionName := range pollForId.Options {
		result[optionName] = []string{}
	}
	for _, vote := range votes {
		optionName := pollForId.Options[vote.VotedFor]
		result[optionName] = append(result[optionName], vote.VoterID)
	}
	return result, nil
}

func (s *DefaultStore) GetPoll(pollId string) (Poll, error) {
	return s.backend.GetPoll(pollId)
}

func (s *DefaultStore) GetVote(voteId string) (Vote, error) {
	return s.backend.GetVote(voteId)
}
