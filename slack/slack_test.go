package slack

import (
	"encoding/json"
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
	"io/ioutil"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Bool("integration", false, "run integration tests")
	flag.Parse()
	os.Exit(m.Run())
}

func TestNewPollMessageSimple(t *testing.T) {
	dat, err := ioutil.ReadFile("examplePollMessage.json")
	if err != nil {
		t.Fatal("Error reading sample file: ", err)
	}
	var expectedPollMessage SlackMessage
	err = json.Unmarshal(dat, &expectedPollMessage)
	if err != nil {
		t.Fatal("Error parsing sample file: ", err)
	}
	poll := poll.Poll{"a712786b-b0c1-45f9-8ba6-816a8b665322", "Test Question", "foobar", []string{"Answer 1", "Answer 2"}}
	actualPollMessage := NewPollMessage(poll, nil)
	if diff := deep.Equal(expectedPollMessage, actualPollMessage); diff != nil {
		t.Logf("Created poll message is not as expected.\nExpected: %v\nActual:%v\n", expectedPollMessage, actualPollMessage)
		t.Log("Diff: ", diff)
		t.Fail()
	}

}

func TestParseSlashCommand(t *testing.T) {
	arguments := []string{
		"",
		"\"",
		" ",
		"a",
		"a      ",
		"a     \"",
		"a b c d",
		"\"a b c d\"",
		"a \"b c d\"",
		"a \"b c\" d",
		"a b\"c d",
	}
	expectedResults := [][]string{
		{},
		{},
		{},
		{"a"},
		{"a"},
		{"a"},
		{"a", "b", "c", "d"},
		{"a b c d"},
		{"a", "b c d"},
		{"a", "b c", "d"},
		{"a", "b\"c", "d"},
	}
	for i := range arguments {
		t.Run(arguments[i], func(t *testing.T) {
			result := ParseSlashCommand(arguments[i])
			expected := expectedResults[i]
			if !compareStringSlices(result, expected) {
				printError(t, result, expected)
			}
		})
	}
}

func printError(t *testing.T, actual, expected interface{}) {
	t.Error("Test failed! Expected: \n", spew.Sdump(expected), "\nbut was:\n", spew.Sdump(actual))
}

func TestCompareStringSlices(t *testing.T) {
	if !compareStringSlices([]string{}, []string{}) {
		t.Fail()
	}
	if !compareStringSlices(nil, nil) {
		t.Fail()
	}
	if compareStringSlices(nil, []string{}) {
		t.Fail()
	}
	if compareStringSlices([]string{}, nil) {
		t.Fail()
	}
	if compareStringSlices([]string{""}, []string{"", ""}) {
		t.Fail()
	}
	if !compareStringSlices([]string{"a", "b", "c"}, []string{"a", "b", "c"}) {
		t.Fail()
	}
	if compareStringSlices([]string{"a", "b", "c"}, []string{"a", "d", "c"}) {
		t.Fail()
	}
}

func compareStringSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (b == nil && a != nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
