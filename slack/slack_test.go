package slack

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"bytes"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
	"markusreschke.name/selfhostedchatpolling/poll"
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
	poll := poll.Poll{ID: "6b57e603-2366-4116-b51d-011837677e33", Question: "Test Question", CreatorID: "foobar", Options: []string{"Answer 1", "Answer 2"}}
	actualPollMessage := NewPollMessage(poll, nil)
	if diff := deep.Equal(expectedPollMessage, actualPollMessage); diff != nil {
		t.Logf("Created poll message is not as expected.\nExpected: %v\nActual:%v\n", expectedPollMessage, actualPollMessage)
		t.Log("Diff: ", diff)
		t.Fail()
	}
}

func TestNewPollMessageAnonymous(t *testing.T) {
	dat, err := ioutil.ReadFile("exampleAnonymousPollMessage.json")
	if err != nil {
		t.Fatal("Error reading sample file: ", err)
	}
	var expectedPollMessage SlackMessage
	err = json.Unmarshal(dat, &expectedPollMessage)
	if err != nil {
		t.Fatal("Error parsing sample file: ", err)
	}
	poll := poll.Poll{ID: "f843f53f-d7d2-4050-a5d7-fd222114038f", Question: "Test Question", CreatorID: "foobar", Options: []string{"Answer 1", "Answer 2"}, Anonymous: true}
	actualPollMessage := NewPollMessage(poll, nil)
	if diff := deep.Equal(expectedPollMessage, actualPollMessage); diff != nil {
		t.Logf("Created poll message is not as expected.\nExpected: %v\nActual:%v\n", expectedPollMessage, actualPollMessage)
		t.Log("Diff: ", diff)
		t.Fail()
	}
}

func TestNewVoteDetailMessage(t *testing.T) {
	expectedText := "• Option1: A, B, C\n• Option2: A, B\n• Option3: \n"
	input := map[string][]string{
		"Option1": {"A", "B", "C"},
		"Option2": {"A", "B"},
		"Option3": {},
	}
	slackMsg := NewVoteDetailMessage(input)
	if diff := deep.Equal(expectedText, slackMsg.Text); diff != nil {
		t.Log("Vote Detail message is not formed as expected!")
		t.Log("Diff:\n", diff)
		t.Fail()
	}
	if slackMsg.ReplaceOriginal != false {
		t.Log("Vote Detail message set to replace original!")
		t.Fail()
	}
	if slackMsg.ResponseType != ResponseTypeEphemeral {
		t.Log("Vote Detail message is not ephemeral!")
		t.Fail()
	}
}

func TestBuildVoteDetailMessageTest(t *testing.T) {
	expectedText := "• Option1: A, B, C\n• Option2: A, B\n• Option3: \n"
	input := map[string][]string{
		"Option1": {"A", "B", "C"},
		"Option2": {"A", "B"},
		"Option3": {},
	}
	var actualText bytes.Buffer
	buildVoteDetailMessageTest(input, &actualText)
	if diff := deep.Equal(expectedText, actualText.String()); diff != nil {
		t.Log("Vote Detail message is not formed as expected!")
		t.Log("Diff:\n", diff)
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
