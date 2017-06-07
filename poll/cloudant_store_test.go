// +build integration
package poll

import (
	"flag"
	"fmt"
	"github.com/IBM-Bluemix/go-cloudant"
	"os"
	"reflect"
	"testing"
)

const (
	testDBName = "test_poll_db"
)

var client *cloudant.Client

func TestMain(m *testing.M) {
	integrationTest := flag.Bool("integration", false, "run integration tests")
	flag.Parse()
	if !*integrationTest {
		fmt.Printf("Skiped test because -integration was not used!\n")
		os.Exit(0)
	}
	var err error
	client, err = cloudant.NewClient(os.Getenv("CLOUDANT_USER"), os.Getenv("CLOUDANT_PW"))
	if err != nil {
		fmt.Printf("Error connecting to cloudant: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func GetCleanStore(client *cloudant.Client) StoreBackend {
	err := client.DeleteDB(testDBName)
	if err != nil {
		fmt.Printf("Error deleting db: %v", err)
	}
	store, err := NewCloudantStore(client, testDBName)
	if err != nil {
		fmt.Printf("Error creating store: %v", err)
		os.Exit(1)
	}
	return store
}

func TestAddingAndRetrievingData(t *testing.T) {
	store := GetCleanStore(client)
	poll := Poll{"1", "q", "creator", []string{"a1", "a2"}}
	err := store.AddPoll(poll)
	if err != nil {
		t.Fatalf("Error creating poll: %v", err)
	}
	vote := Vote{"1", "voter", "1", "a1"}
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

func TestGettingVotesForPoll(t *testing.T) {
	store := GetCleanStore(client)
	votes := []Vote{
		{"1", "voter", "1", "a1"},
		{"2", "voter2", "1", "a1"},
		{"3", "voter3", "1", "a3"},
	}
	poll := Poll{"1", "q", "creator", []string{"a1", "a2", "a3"}}
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
