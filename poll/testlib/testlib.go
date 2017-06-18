package testlib

import (
	. "markusreschke.name/selfhostedchatpolling/poll"
	"reflect"
	"testing"
)

func TestAddingAndRetrievingData(t *testing.T, store StoreBackend) {
	poll := Poll{"1", "q", "creator", []string{"a1", "a2"}}
	err := store.AddPoll(poll)
	if err != nil {
		t.Fatalf("Error creating poll: %v", err)
	}
	vote := Vote{"1", "voter", "1", 0}
	err = store.AddVote(vote)
	if err != nil {
		t.Fatalf("Error creating vote: %v", err)
	}
	pollFromStore, err := store.GetPoll("1")
	if !reflect.DeepEqual(poll, pollFromStore) {
		t.Fatalf("Expected %v but got %v", poll, pollFromStore)
	}
	voteFromStore, err := store.GetVote("1")
	if !reflect.DeepEqual(vote, voteFromStore) {
		t.Fatalf("Expected %v but got %v", vote, voteFromStore)
	}
}

func TestPollHasVoteFromVoter(t *testing.T, store StoreBackend) {
	poll := Poll{"1", "q", "creator", []string{"a1", "a2"}}
	voterID := "voter"
	isVotePresent, err := store.PollHasVoteFromVoter(poll.ID, voterID)
	if err != nil || isVotePresent {
		if err != nil {
			t.Fatalf("Error looking up vote: %v", err)
		} else {
			t.Fatalf("Vote found for Voter %s in empty store!", voterID)
		}
	}
	vote := Vote{"1", voterID, "1", 0}
	err = store.AddVote(vote)
	if err != nil {
		t.Fatalf("Error creating vote: %v", err)
	}
	isVotePresent, err = store.PollHasVoteFromVoter(poll.ID, voterID)
	if err != nil || !isVotePresent {
		if err != nil {
			t.Fatalf("Error looking up vote: %v", err)
		} else {
			t.Fatalf("Vote not found for Voter %s after adding it!", voterID)
		}
	}
}

func TestGettingVotesForPoll(t *testing.T, store StoreBackend) {
	poll := Poll{"1", "q", "creator", []string{"a1", "a2", "a3"}}
	votes := []Vote{
		{"1", "voter", "1", 0},
		{"2", "voter2", "1", 0},
		{"3", "voter3", "1", 2},
	}
	store.AddPoll(poll)
	for _, vote := range votes {
		store.AddVote(vote)
	}
	result, err := store.GetVotesForPoll(poll.ID)
	if err != nil {
		t.Fatalf("Error while fetching votes for Poll %s: %v", poll.ID, err)
	}
	compareVotes(t, votes, result)
}

func compareVotes(t *testing.T, expected, actual []Vote) {
expectedLoop:
	for _, expectedVote := range expected {
		for _, actualVote := range actual {
			if expectedVote.ID == actualVote.ID {
				if expectedVote == actualVote {
					continue expectedLoop
				} else {
					t.Fatalf("Expected %v but got %v", expectedVote, actualVote)
				}
			}
		}
		t.Fatalf("No matching vote found for ID %s", expectedVote.ID)
	}
}
