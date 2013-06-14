package markov

import (
	"fmt"
	"testing"
)

var splitTests = []struct {
	in       string
	expected []string
}{
	{"one", []string{"one"}},
	{"one two", []string{"one", "two"}},
	{"I am the one and only.", []string{"I", "am", "the", "one", "and", "only."}},
}

func verify(t *testing.T, testnum int, testcase, input string, output, expected []string) {
	outputString := fmt.Sprintf("%v", output)
	expectedString := fmt.Sprintf("%v", expected)
	if outputString != expectedString {
		t.Errorf("%d. %s with input = %s: output %s != %s", testnum, testcase, input, output, expected)
	}
}

func TestSplit(t *testing.T) {
	chain := NewChain(1)
	for i, tt := range splitTests {
		output := chain.splitString(tt.in)
		verify(t, i, "TestSplit", tt.in, output, tt.expected)
	}
}

func verifyMap(t *testing.T, testnum int, testcase string, input map[string][]string, key, expected string) {
	v, ok := input[key]
	if !ok {
		t.Errorf("%d. %s no value found for key \"%s\"", testnum, testcase, key)
	} else if fmt.Sprintf("%v", v) != expected {
		t.Errorf("%d. %s with key = %s: value %s != %s", testnum, testcase, key, v, expected)
	}
}

func TestCreateChain(t *testing.T) {
	chain := NewChain(1)
	if chain == nil {
		t.Errorf("Chain not created")
	} else if chain.grams == nil {
		t.Errorf("grams not created")
	} else if chain.order != 1 {
		t.Errorf("order incorrect. value %s != %s", chain.order, 1)
	}
}

func TestPopulateTableOrder1(t *testing.T) {
	c := NewChain(1)
	c.Populate("one two two three")
	verifyMap(t, 1, "TestPopulateTableOrder1", c.grams, "", "[one]")
	verifyMap(t, 2, "TestPopulateTableOrder1", c.grams, "one", "[two]")
	verifyMap(t, 3, "TestPopulateTableOrder1", c.grams, "two", "[two three]")
	verifyMap(t, 4, "TestPopulateTableOrder1", c.grams, "three", "[TERMINATORTERMINATOR]")
}

func TestPopulateTableOrder2(t *testing.T) {
	c := NewChain(2)
	c.Populate("one two two three")
	verifyMap(t, 1, "TestPopulateTableOrder1", c.grams, " ", "[one]")
	verifyMap(t, 2, "TestPopulateTableOrder1", c.grams, " one", "[two]")
	verifyMap(t, 3, "TestPopulateTableOrder1", c.grams, "one two", "[two]")
	verifyMap(t, 4, "TestPopulateTableOrder1", c.grams, "two two", "[three]")
	verifyMap(t, 5, "TestPopulateTableOrder1", c.grams, "two three", "[TERMINATORTERMINATOR]")
}

func TestPopulateTableOrder3(t *testing.T) {
	c := NewChain(3)
	c.Populate("one two two three")
	verifyMap(t, 1, "TestPopulateTableOrder1", c.grams, "  ", "[one]")
	verifyMap(t, 2, "TestPopulateTableOrder1", c.grams, "  one", "[two]")
	verifyMap(t, 3, "TestPopulateTableOrder1", c.grams, " one two", "[two]")
	verifyMap(t, 4, "TestPopulateTableOrder1", c.grams, "one two two", "[three]")
	verifyMap(t, 5, "TestPopulateTableOrder1", c.grams, "two two three", "[TERMINATORTERMINATOR]")
}

func TestPopulateTableOrder1MultipleChildren(t *testing.T) {
	c := NewChain(1)
	c.Populate("One two. One two. One two three.")
	verifyMap(t, 1, "TestPopulateTableOrder1", c.grams, "", "[One]")
	verifyMap(t, 2, "TestPopulateTableOrder1", c.grams, "One", "[two. two. two]")
	verifyMap(t, 3, "TestPopulateTableOrder1", c.grams, "two.", "[One One]")
	verifyMap(t, 4, "TestPopulateTableOrder1", c.grams, "two", "[three.]")
	verifyMap(t, 5, "TestPopulateTableOrder1", c.grams, "three.", "[TERMINATORTERMINATOR]")
}

/*
func TestGenerate(t * testing.T) {
	c := NewChain(1)
	c.Populate("
}
*/
