package slack

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Bool("integration", false, "run integration tests")
	flag.Parse()
	os.Exit(m.Run())
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
