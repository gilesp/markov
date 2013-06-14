package markov

import (
	"math/rand"
	"strings"
	"time"
)

/*
 * A Splitter to aid in dub fiction remixing.
 * See http://www.metamorphiction.com/index.php/the-ghost-on-the-b-side-remixing-narrative for more details
 *
 * This splitter works at the word level, to produce a more cut up series
 * of outputs than a sentence-oriented splitter.
 */
type DubSplitter struct {
	phraseLength int
}

func (d *DubSplitter) Split(text string) []string {

	phrases := []string{}
	current := []string{}
	words := strings.Fields(text)
	for _, word := range words {
		if len(current) < d.phraseLength {
			current = append(current, word)
		} else {
			phrases = append(phrases, d.makePhrase(current))
			current = []string{}
		}
	}

	//catch any scraps left over
	if len(current) > 0 {
		phrases = append(phrases, d.makePhrase(current))
	}
	return phrases
}

func (d *DubSplitter) makePhrase(words []string) string {
	return strings.Join(words, " ")
}

func NewDubSplitter() *DubSplitter {
	d := &DubSplitter{}
	//	d.phraseLength = 4 //TODO: randomise this to between 2 and 5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d.phraseLength = r.Intn(5) + 2
	return d
}
