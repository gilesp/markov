//A Very rough and ready go port of the University of Illinois' Cognitive Computation Groups's sentence segmentation tool
//The original can be found here: http://cogcomp.cs.illinois.edu/page/tools_view/2
//This implementation is quick and dirty and not terribly well written - it really needs tidying up and to make better use of go's idioms.
package markov

import (
	"regexp"
	"strings"
)

type IllinoisSplitter struct{}

var honorifics = []string{
	"A.",
	"Adj.",
	"Adm.",
	"Adv.",
	"Asst.",
	"B.",
	"Bart.",
	"Bldg.",
	"Brig.",
	"Bros.",
	"C.",
	"Capt.",
	"Cmdr.",
	"Col.",
	"Comdr.",
	"Con.",
	"Cpl.",
	"D.",
	"DR.",
	"Dr.",
	"E.",
	"Ens.",
	"F.",
	"G.",
	"Gen.",
	"Gov.",
	"H.",
	"Hon.",
	"Hosp.",
	"I.",
	"Insp.",
	"J.",
	"K.",
	"L.",
	"Lt.",
	"M.",
	"M.",
	"MM.",
	"MR.",
	"MRS.",
	"MS.",
	"Maj.",
	"Messrs.",
	"Mlle.",
	"Mme.",
	"Mr.",
	"Mrs.",
	"Ms.",
	"Msgr.",
	"N.",
	"O.",
	"Op.",
	"Ord.",
	"P.",
	"Pfc.",
	"Ph.",
	"Prof.",
	"Pvt.",
	"Q.",
	"R.",
	"Rep.",
	"Reps.",
	"Res.",
	"Rev.",
	"Rt.",
	"S.",
	"Sen.",
	"Sens.",
	"Sfc.",
	"Sgt.",
	"Sr.",
	"St.",
	"Supt.",
	"Surg.",
	"T.",
	"U.",
	"V.",
	"W.",
	"X.",
	"Y.",
	"Z.",
	"v.",
	"vs.",
}

func (is *IllinoisSplitter) Split(text string) []string {
	words := strings.Fields(text)
	sentences := []string{}
	sentence := ""

	for i, word := range words {

		//check the existence of a candidate
		periodPos := strings.LastIndex(word, ".")
		questionPos := strings.LastIndex(word, "?")
		exclaimPos := strings.LastIndex(word, "!")

		//which is the latest in the word?
		pos := periodPos
		candidate := "."
		if questionPos > pos {
			pos = questionPos
			candidate = "?"
		}
		if exclaimPos > pos {
			pos = exclaimPos
			candidate = "!"
		}

		if pos > -1 {
			//check for previous word(s)
			var wm1 string
			var wm2 string
			if i > 0 {
				wm1 = words[i-1]

				if i > 1 {
					wm2 = words[i-2]
				}
			}

			//check for next word(s)
			var wp1 string
			var wp2 string
			if i < (len(words) - 1) {
				wp1 = words[i+1]
				if i < (len(words) - 2) {
					wp2 = words[i+2]
				}
			}

			//define prefix
			prefix := "sp"
			if pos > 0 {
				prefix = word[0:pos]
			}
			//define suffix
			suffix := "sp"
			if pos != (len(word) - 1) {
				suffix = word[pos+1 : len(word)-1]
			}

			//append word to sentence
			sentence = is.appendString(sentence, word)

			if is.boundary(candidate, wm2, wm1, prefix, suffix, wp1, wp2) {
				sentences = append(sentences, sentence)
				sentence = ""
			}
		} else {
			sentence = is.appendString(sentence, word)
		}

	}
	return sentences
}

func (i *IllinoisSplitter) appendString(sentence, word string) string {
	if sentence == "" {
		sentence = word
	} else {
		sentence = sentence + " " + word
	}
	return sentence
}

func (i *IllinoisSplitter) startsWithCapital(word string) bool {
	if len(word) > 0 {
		matched, _ := regexp.MatchString("[A-Z]", word[:1])
		return matched
	} else {
		return false
	}
}

func (i *IllinoisSplitter) boundary(candidate, prev2Word, prevWord, prefix, suffix, nextWord, next2Word string) bool {
	if candidate == "?" || candidate == "!" {
		if nextWord == "" && next2Word == "" {
			return true
		}

		if suffix == "sp" {
			if i.startsWithCapital(nextWord) {
				return true
			}

			if i.startsWithQuote(nextWord) {
				return true
			}

			if nextWord == "--" && i.startsWithCapital(next2Word) {
				return true
			}
			if nextWord == "." {
				return true
			}
		}
		if i.isRightEnd(suffix) && i.isLeftStart(nextWord) {
			return true
		} else {
			return false
		}

	} else {
		if nextWord == "" && next2Word == "" {
			return true
		}
		if suffix == "sp" {
			if i.startsWithQuote(nextWord) {
				return true
			}

			if i.startsWithLeftParen(nextWord) {
				return true
			}

			if next2Word == "--" {
				return false
			}

			if i.isRightParen(nextWord) {
				return true
			}

			if candidate == "." && i.endsWithRightParen(nextWord) && i.startsWithCapital(next2Word) {
				return true
			}

			if prefix == "sp" && nextWord == "." {
				return false
			}

			if nextWord == "." {
				return true
			}

			if nextWord == "--" {
				if i.startsWithCapital(next2Word) {
					if i.endsInQuote(prefix) {
						return true
					}
				} else if i.startsWithQuote(next2Word) {
					return true
				}
			}

			if i.startsWithCapital(nextWord) && (prefix == "p.m" || prefix == "a.m") {
				return false
			}

			if i.startsWithCapital(nextWord) && (i.isHonorific(prefix+".") || i.startsWithQuote(prefix)) {
				return false
			}

			if i.startsWithCapital(nextWord) && i.isTerminal(prefix) {
				return true
			}

			if i.startsWithCapital(nextWord) {
				return true
			}
		}

		if i.isRightEnd(suffix) && i.isLeftStart(nextWord) {
			return true
		}
	}
	return false
}

func (i *IllinoisSplitter) isHonorific(word string) bool {
	for _, honorific := range honorifics {
		if word == honorific {
			return true
		}
	}
	return false
}

func (i *IllinoisSplitter) isTerminal(word string) bool {
	terminals := []string{"Esq", "Jr", "Sr", "M.D", "Phd"}
	for _, terminal := range terminals {
		if word == terminal {
			return true
		}
	}
	return false
}

func (i *IllinoisSplitter) endsInQuote(word string) bool {
	runes := []rune(word)
	length := len(runes) - 1
	lastChar := runes[length]
	return i.isQuote(string(lastChar))
}

func (i *IllinoisSplitter) startsWithQuote(word string) bool {
	runes := []rune(word)
	firstChar := runes[0]
	return i.isQuote(string(firstChar))
}

func (i *IllinoisSplitter) isQuote(char string) bool {
	switch string(char) {
	case "''":
		fallthrough
	case "'":
		fallthrough
	case "'''":
		fallthrough
	case "\"":
		fallthrough
	case "'\"":
		return true
	default:
		return false
	}
}

func (i *IllinoisSplitter) startsWithLeftParen(word string) bool {
	runes := []rune(word)
	firstChar := string(runes[0])
	return firstChar == "(" || firstChar == "{" || firstChar == "[" || firstChar == "<"
}

func (i *IllinoisSplitter) endsWithRightParen(word string) bool {
	runes := []rune(word)
	return i.isRightParen(string(runes[len(runes)-1]))
}

func (i *IllinoisSplitter) startsWithLeftQuote(word string) bool {
	runes := []rune(word)
	firstChar := string(runes[0])
	return firstChar == "`" || firstChar == "\"" || firstChar == "\"`"
}

func (i *IllinoisSplitter) isRightEnd(word string) bool {
	return i.isRightParen(word) || i.isRightQuote(word)
}

func (i *IllinoisSplitter) isLeftStart(word string) bool {
	return i.startsWithLeftQuote(word) || i.startsWithLeftParen(word)
}

func (i *IllinoisSplitter) isRightParen(word string) bool {
	return word == ")" || word == "}" || word == "]" || word == ">"
}

func (i *IllinoisSplitter) isRightQuote(word string) bool {
	return word == "'" || word == "''" || word == "'''" || word == "\"" || word == "'\""
}

func NewIllinoisSplitter() *IllinoisSplitter {
	s := &IllinoisSplitter{}
	return s
}
