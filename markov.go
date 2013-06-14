package markov

import (
	"math/rand"
	"strings"
	"time"
)

// key definitions
type key []string

func (k key) String() string {
	return strings.Join(k, " ")
}

func (k key) shift(word string) {
	copy(k, k[1:])
	k[len(k)-1] = word
}

const terminator = "TERMINATORTERMINATOR"

//Chain definitions
type Chain struct {
	grams    map[string][]string
	order    int
	splitter Splitter
}

func (c *Chain) splitString(text string) []string {
	return strings.Fields(text)
}

func (c *Chain) Size() int {
	return len(c.grams)
}

func (c *Chain) Populate(text string) {
	var sentences []string
	if c.splitter != nil {
		sentences = c.splitter.Split(text)
	} else {
		sentences = append(sentences, text)
	}
	for _, sentence := range sentences {
		words := c.splitString(sentence)
		words = append(words, terminator)
		key := make(key, c.order)
		for _, word := range words {
			c.grams[key.String()] = append(c.grams[key.String()], word)
			key.shift(word)
		}
	}
}

func (c *Chain) generateWithKey(key key, maxLength int) string {

	words := []string{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sentence := []string{}

	for i := 0; i < maxLength; i++ {
		choices := c.grams[key.String()]
		if len(choices) == 0 {
			//we have reached the end of a chain
			//reset the key and carry on.
			for j := 0; j < c.order; j++ {
				key.shift("")
			}
			choices = c.grams[key.String()]
		}
		next := choices[r.Intn(len(choices))]

		//if next is a terminator word
		//append sentence to words.
		//this is to aovid half-finished sentences appearing in the output.
		if next == terminator {
			words = append(words, strings.Join(sentence, " "))
			sentence = []string{}
		} else {
			sentence = append(sentence, next)
		}
		key.shift(next)
	}

	return strings.Join(words, " ")
}

func (c *Chain) GenerateResponse(stub string, maxLength int) (response string, err string) {
	words := strings.Fields(stub)
	if len(words) < c.order {
		err = "Stub too short."
		return
	}
	key := make(key, c.order)
	for i := (len(words) - c.order); i < len(words); i++ {
		key.shift(words[i])
	}
	response = c.generateWithKey(key, maxLength)
	return
}

func (c *Chain) Generate(maxLength int) string {
	key := make(key, c.order)

	return c.generateWithKey(key, maxLength)
}

// Public utility method to create a new Chain instance
func NewChain(order int) *Chain {
	c := &Chain{
		grams: make(map[string][]string),
		order: order,
	}
	return c
}

func NewChainWithSplitter(order int, splitter Splitter) *Chain {
	c := &Chain{
		grams:    make(map[string][]string),
		order:    order,
		splitter: splitter,
	}
	return c
}
