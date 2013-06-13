package markov

import (
	"fmt"
	"testing"
)

var sentenceTests = []struct {
	in       string
	expected []string
}{
	{"This is one sentence.", []string{"This is one sentence."}},
	{"This is sentence one. This is sentence two.", []string{"This is sentence one.", "This is sentence two."}},
	{"Dr. this is one sentence.", []string{"Dr. this is one sentence."}},
	{"Prof. Marshall is very old. His wife, Mrs. Marshall is very young.", []string{"Prof. Marshall is very old.", "His wife, Mrs. Marshall is very young."}},
}

func verifyStringSlice(t *testing.T, testnum int, testcase, input string, output, expected []string) {
	outputString := fmt.Sprintf("%v", output)
	expectedString := fmt.Sprintf("%v", expected)
	if outputString != expectedString {
		t.Errorf("%d. %s with input = %s: output %s != %s", testnum, testcase, input, output, expected)
	}
}

func verifyBool(t *testing.T, testnum int, testcase, input string, output, expected bool) {
	if output != expected {
		t.Errorf("%d. %s with input = \"%s\": expected %t, found %t", testnum, testcase, input, expected, output)
	}
}

func TestNaiveSplit(t *testing.T) {
	s := NewNaiveSplitter()
	for i, st := range sentenceTests {
		output := s.Split(st.in)
		verifyStringSlice(t, i, "TestSplit", st.in, output, st.expected)
	}
}

var terminatorTests = []struct {
	in       string
	expected bool
}{
	{"Help!", true},
	{"Help.", true},
	{"Help?", true},
	{"Help", false},
	{"Mr.", false},
	{"Mrs.", false},
	{"Ms.", false},
	{"Miss.", false},
	{"Messrs.", false},
	{"Jr.", false},
	{"Snr.", false},
	{"Dr.", false},
	{"Prof.", false},
}

func TestIsTerminator(t *testing.T) {
	s := NewNaiveSplitter()
	for i, tt := range terminatorTests {
		verifyBool(t, i, "TestIsTerminator", tt.in, s.IsTerminator(tt.in), tt.expected)
	}
}
