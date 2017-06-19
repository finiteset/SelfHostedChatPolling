package memstore

import (
	"flag"
	"markusreschke.name/selfhostedchatpolling/poll"
	"markusreschke.name/selfhostedchatpolling/poll/testlib"
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

func TestAllCasesInTestLib(t *testing.T) {
	testlib.RunTests(t, getCleanStore)
}
