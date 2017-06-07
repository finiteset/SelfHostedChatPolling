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
	if !*integrationTest {
		fmt.Printf("Skiped test because -interation was not used!\n")
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

func GetCleanStore(client *cloudant.Client) Store {
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

func TestGettingCount(t *testing.T) {
	store := GetCleanStore(client)
	poll := Poll{"1", "q", "creator", []string{"a1", "a2", "a3"}}
	store.AddPoll(poll)
	vote := Vote{"1", "voter", "1", "a1"}
	store.AddVote(vote)
	vote = Vote{"2", "voter2", "1", "a1"}
	store.AddVote(vote)
	vote = Vote{"3", "voter3", "1", "a3"}
	store.AddVote(vote)
	vote = Vote{"4", "voter4", "1", "a1"}
	store.AddVote(vote)
	vote = Vote{"5", "voter5", "1", "a3"}
	store.AddVote(vote)
	vote = Vote{"6", "voter6", "1", "a1"}
	store.AddVote(vote)
	result, err := store.GetResult("1")
	if err != nil || result["a1"] != 4 || result["a2"] != 0 || result["a3"] != 2 {
		t.Fatalf("Counts do not match. Error: %v", err)
	}
}
