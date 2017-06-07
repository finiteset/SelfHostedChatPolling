package memstore

import (
	"flag"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"markusreschke.name/selfhostedsimplepolling/poll/testlib"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Bool("integration", false, "run integration tests")
	flag.Parse()
	os.Exit(m.Run())
}

func getCleanStore() poll.StoreBackend {
	return NewInMemoryStoreBackend()
}

func TestAddingAndRetrievingData(t *testing.T) {
	testlib.TestAddingAndRetrievingData(t, getCleanStore())
}

func TestGettingVotesForPoll(t *testing.T) {
	testlib.TestGettingVotesForPoll(t, getCleanStore())
}
