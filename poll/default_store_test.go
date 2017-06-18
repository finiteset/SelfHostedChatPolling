package poll_test

import (
	"flag"
	"markusreschke.name/selfhostedchatpolling/poll"
	"markusreschke.name/selfhostedchatpolling/poll/memstore"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Bool("integration", false, "run integration tests")
	flag.Parse()
	os.Exit(m.Run())
}

func TestAddingAndRetrievingData(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{"1", "q", "creator", []string{"a1", "a2"}}
	err := store.AddPoll(testPoll)
	if err != nil {
		t.Fatalf("Error creating poll: %v", err)
	}
	vote := poll.Vote{"1", "voter", "1", 0}
	err = store.AddVote(vote)
	if err != nil {
		t.Fatalf("Error creating vote: %v", err)
	}
	pollFromStore, err := store.GetPoll("1")
	if !reflect.DeepEqual(testPoll, pollFromStore) {
		t.Fatalf("Expected %v but got %v", testPoll, pollFromStore)
	}
	voteFromStore, err := store.GetVote("1")
	if !(vote == voteFromStore) {
		t.Fatalf("Expected %v but got %v", vote, voteFromStore)
	}

	// Test if store allows for invalid voting
	voteInvalidDouble := poll.Vote{"1", "voter", "1", 1}
	err = store.AddVote(voteInvalidDouble)
	if err == nil {
		t.Error("Store allowed double voting by same voter")
	}
	voteInvalidChoice := poll.Vote{"1", "voter2", "1", 2}
	err = store.AddVote(voteInvalidChoice)
	if err == nil {
		t.Error("Store allowed voting for invalid choice")
	}
}

func TestGettingCount(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{"1", "q", "creator", []string{"a1", "a2", "a3"}}
	store.AddPoll(testPoll)
	vote := poll.Vote{"1", "voter", "1", 0}
	store.AddVote(vote)
	vote = poll.Vote{"2", "voter2", "1", 0}
	store.AddVote(vote)
	vote = poll.Vote{"3", "voter3", "1", 2}
	store.AddVote(vote)
	vote = poll.Vote{"4", "voter4", "1", 0}
	store.AddVote(vote)
	vote = poll.Vote{"5", "voter5", "1", 2}
	store.AddVote(vote)
	vote = poll.Vote{"6", "voter6", "1", 0}
	store.AddVote(vote)
	result, err := store.GetResult("1")
	if err != nil || result[0] != 4 || result[1] != 0 || result[2] != 2 {
		t.Fatalf("Counts do not match. Error: %v", err)
	}
}
