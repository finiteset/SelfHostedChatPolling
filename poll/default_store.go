package poll

import (
	"errors"
	"fmt"
)

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
	doubleVote, err := s.backend.PollHasVoteFromVoter(v.PollID, v.VoterID)
	if err != nil {
		return err
	}
	if doubleVote {
		return errors.New(fmt.Sprintf("Voter %s has already voted on Poll %s", v.VoterID, v.PollID))
	}
	isValidChoice, err := s.votedForValidOption(v)
	if err != nil {
		return err
	}
	if !isValidChoice {
		return errors.New(fmt.Sprintf("Voter %s voted for invalid choice %s", v.VoterID, v.VotedFor))
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

func (s *DefaultStore) GetPoll(pollId string) (Poll, error) {
	return s.backend.GetPoll(pollId)
}

func (s *DefaultStore) GetVote(voteId string) (Vote, error) {
	return s.backend.GetVote(voteId)
}
