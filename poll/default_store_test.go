package poll_test

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
	"markusreschke.name/selfhostedchatpolling/poll"
	"markusreschke.name/selfhostedchatpolling/poll/memstore"
)

func TestMain(m *testing.M) {
	flag.Bool("integration", false, "run integration tests")
	flag.Parse()
	os.Exit(m.Run())
}

func TestAddingAndRetrievingData(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{ID: "1", Question: "q", CreatorID: "creator", Options: []string{"a1", "a2", "a3"}}
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
	voteInvalidChoice := poll.Vote{"1", "voter2", "1", 3}
	err = store.AddVote(voteInvalidChoice)
	if err == nil {
		t.Error("Store allowed voting for invalid choice")
	}
}

func TestChangingVotes(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{ID: "1", Question: "q", CreatorID: "creator", Options: []string{"a1", "a2", "a3"}}
	store.AddPoll(testPoll)
	vote := poll.Vote{"1", "voter", "1", 0}
	err := store.AddVote(vote)
	if err != nil {
		t.Log("Error storing first vote: ", err)
		t.Fail()
	}
	vote = poll.Vote{"1", "voter", "1", 1}
	err = store.AddVote(vote)
	if err != nil {
		t.Log("Error storing changed vote: ", err)
		t.Fail()
	}
	result, err := store.GetResult("1")
	if err != nil || result[0] != 0 || result[1] != 1 || result[2] != 0 {
		t.Fatal("Counts do not match. Error: ", err)
	}
}

func TestGettingCount(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{ID: "1", Question: "q", CreatorID: "creator", Options: []string{"a1", "a2", "a3"}}
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
		t.Fatal("Counts do not match. Error: ", err)
	}
}

func TestGetVoteDetails(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{ID: "1", Question: "q", CreatorID: "creator", Options: []string{"a1", "a2", "a3"}}
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

	expectedResult := map[string][]string{
		"a1": []string{"voter", "voter2", "voter4", "voter6"},
		"a2": []string{},
		"a3": []string{"voter3", "voter5"},
	}

	result, err := store.GetVoteDetails(testPoll.ID)
	if err != nil {
		t.Fatal("Error getting vote details: ", err)
	}
	if diff := deep.Equal(expectedResult, result); diff != nil {
		spew.Dump(result)
		t.Log("Calculated vote details don't match the expected details", diff)
		t.Fail()
	}
}

func TestGetVoteDetailsAnonymous(t *testing.T) {
	store := poll.NewDefaultStore(memstore.NewInMemoryStoreBackend())
	testPoll := poll.Poll{ID: "1", Question: "q", CreatorID: "creator", Options: []string{"a1", "a2", "a3"}, Anonymous: true}
	store.AddPoll(testPoll)

	_, err := store.GetVoteDetails(testPoll.ID)
	if err == nil {
		t.Fatal("Was able to get vote details for anonymous poll!")
	}
	if err != poll.ErrNoDetailsForAnonymousPoll {
		t.Fatal("Unexpected error was returned!: ", err)
	}
}
