package cloudantstore

import (
	"strings"

	"github.com/IBM-Bluemix/go-cloudant"
	"markusreschke.name/selfhostedchatpolling/poll"
)

type CloudantStore struct {
	db *cloudant.DB
}

const (
	pollPrefix    = "poll_"
	votePrefix    = "vote_"
	deleteRetries = 3
)

func buildCloudantVoteId(voteId string) string {
	return votePrefix + voteId
}

func NewCloudantStoreBackend(client *cloudant.Client, dbName string) (poll.StoreBackend, error) {
	db, err := client.CreateDB(dbName)
	if err != nil {
		db, err = client.EnsureDB(dbName)
		if err != nil {
			return nil, err
		}
	}
	return &CloudantStore{db}, nil
}

func (s *CloudantStore) AddPoll(p poll.Poll) error {
	p.ID = pollPrefix + p.ID
	_, _, err := s.db.CreateDocument(p)
	return err
}

func (s *CloudantStore) AddVote(v poll.Vote) error {
	v.ID = buildCloudantVoteId(v.ID)
	_, _, err := s.db.CreateDocument(v)
	return err
}

func rebuildVotesFromSearchResult(votes []interface{}) ([]poll.Vote, error) {
	result := []poll.Vote{}
	for _, rawVote := range votes {
		voteMap := rawVote.(map[string]interface{})
		vote, err := rebuildVoteFromMap(voteMap)
		if err != nil {
			return nil, err
		}
		result = append(result, vote)
	}
	return result, nil
}

func rebuildVoteFromMap(voteMap map[string]interface{}) (poll.Vote, error) {
	return poll.Vote{
		strings.TrimPrefix(voteMap["_id"].(string), votePrefix),
		voteMap["VoterID"].(string),
		voteMap["PollID"].(string),
		int(voteMap["VotedFor"].(float64)),
	}, nil
}

func (s *CloudantStore) GetVotesForPoll(pollId string) ([]poll.Vote, error) {
	query := cloudant.Query{}
	query.Selector = make(map[string]interface{})
	query.Selector["PollID"] = pollId
	votes, err := s.db.SearchDocument(query)
	if err != nil {
		return nil, err
	}
	return rebuildVotesFromSearchResult(votes)
}

func (s *CloudantStore) GetPoll(pollId string) (poll.Poll, error) {
	var poll poll.Poll
	err := s.db.GetDocument(pollPrefix+pollId, &poll, nil)
	poll.ID = strings.Replace(poll.ID, pollPrefix, "", 1)
	return poll, err
}

func (s *CloudantStore) GetVote(voteId string) (poll.Vote, error) {
	var vote poll.Vote
	err := s.db.GetDocument(buildCloudantVoteId(voteId), &vote, nil)
	vote.ID = strings.Replace(vote.ID, votePrefix, "", 1)
	return vote, err
}

func (s *CloudantStore) PollHasVoteFromVoter(pollID, voterID string) (bool, poll.Vote, error) {
	query := cloudant.Query{}
	query.Selector = make(map[string]interface{})
	query.Selector["PollID"] = pollID
	query.Selector["VoterID"] = voterID
	votes, err := s.db.SearchDocument(query)
	if err != nil {
		return false, poll.Vote{}, err
	}
	if votes == nil || len(votes) == 0 {
		return false, poll.Vote{}, nil
	} else {
		foundVote, err := rebuildVoteFromMap(votes[0].(map[string]interface{}))
		if err != nil {
			return true, poll.Vote{}, err
		}
		return true, foundVote, nil
	}
}

func (s *CloudantStore) RemoveVote(voteId string) error {
	cloudantVoteId := buildCloudantVoteId(voteId)
	var err error = nil
	for i := 1; i <= deleteRetries; i += 1 {
		rev, err := s.db.GetDocumentRev(cloudantVoteId)
		if err != nil {
			continue
		}
		_, err = s.db.DeleteDocument(cloudantVoteId, rev)
		if err == nil {
			break
		}
	}
	return err
}
