// +build integration
package cloudantstore

import (
	"flag"
	"fmt"
	"github.com/IBM-Bluemix/go-cloudant"
	"markusreschke.name/selfhostedchatpolling/poll"
	"markusreschke.name/selfhostedchatpolling/poll/testlib"
	"os"
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
		fmt.Printf("Skipping cloudant test because -integration was not used!\n")
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

func getCleanStore(client *cloudant.Client) poll.StoreBackend {
	err := client.DeleteDB(testDBName)
	if err != nil {
		fmt.Printf("Error deleting db: %v", err)
	}
	store, err := NewCloudantStoreBackend(client, testDBName)
	if err != nil {
		fmt.Printf("Error creating store: %v", err)
		os.Exit(1)
	}
	return store
}

func TestAddingAndRetrievingData(t *testing.T) {
	testlib.TestAddingAndRetrievingData(t, getCleanStore(client))
}

func TestGettingVotesForPoll(t *testing.T) {
	testlib.TestGettingVotesForPoll(t, getCleanStore(client))
}

func TestPollHasVoteFromVoter(t *testing.T) {
	testlib.TestPollHasVoteFromVoter(t, getCleanStore(client))
}
