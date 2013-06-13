package markov

import (
	"strings"
)

type Splitter interface {
	Split(text string) []string
}

type NaiveSplitter struct{}

func (n *NaiveSplitter) Split(text string) []string {

	sentences := []string{}
	current := []string{}
	words := strings.Fields(text)
	for _, word := range words {
		if n.isTerminator(word) {
			current = append(current, word)
			sentences = append(sentences, strings.Join(current, " "))
			current = []string{}
		} else {
			current = append(current, word)
		}
	}
	return sentences
}

var abbrevs = []string{"Mr.",
	"Mrs.",
	"Ms.",
	"Miss.",
	"Messrs.",
	"Jr.",
	"Snr.",
	"Dr.",
	"Prof."}

func (n *NaiveSplitter) isTerminator(word string) bool {

	for _, abbrev := range abbrevs {
		if word == abbrev {
			return false
		}
	}

	runes := []rune(word)
	length := len(runes) - 1
	lastChar := runes[length]
	switch string(lastChar) {
	case "?":
		fallthrough
	case "!":
		fallthrough
	case ".":
		return true
	default:
		return false
	}
}

func NewNaiveSplitter() *NaiveSplitter {
	s := &NaiveSplitter{}
	return s
}
