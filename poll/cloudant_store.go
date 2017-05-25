package poll

import (
	"github.com/IBM-Bluemix/go-cloudant"
	"strings"
)

type CloudantStore struct {
	db *cloudant.DB
}

const (
	pollPrefix = "poll_"
	votePrefix = "vote_"
)

func NewCloudantStore(client *cloudant.Client, dbName string) (Store, error) {
	db, err := client.CreateDB(dbName)
	if err != nil {
		return nil, err
	}
	return &CloudantStore{db}, nil
}

func (s *CloudantStore) AddPoll(p Poll) error {
	p.ID = pollPrefix + p.ID
	_, _, err := s.db.CreateDocument(p)
	return err
}

func (s *CloudantStore) AddVote(v Vote) error {
	v.ID = votePrefix + v.ID
	_, _, err := s.db.CreateDocument(v)
	return err
}

func (s *CloudantStore) GetResult(pollId string) (map[string]uint64, error) {
	result := make(map[string]uint64)
	query := cloudant.Query{}
	query.Selector = make(map[string]interface{})
	query.Selector["PollID"] = pollId
	votes, err := s.db.SearchDocument(query)
	if err != nil {
		return nil, err
	}
	for _, voteRaw := range votes {
		vote := voteRaw.(map[string]interface{})
		result[vote["VotedFor"].(string)]++
	}
	return result, nil
}

func (s *CloudantStore) GetPoll(pollId string) (Poll, error) {
	var poll Poll
	err := s.db.GetDocument(pollPrefix+pollId, &poll, nil)
	poll.ID = strings.Replace(poll.ID, pollPrefix, "", 1)
	return poll, err
}

func (s *CloudantStore) GetVote(voteId string) (Vote, error) {
	var vote Vote
	err := s.db.GetDocument(votePrefix+voteId, &vote, nil)
	vote.ID = strings.Replace(vote.ID, votePrefix, "", 1)
	return vote, err
}
