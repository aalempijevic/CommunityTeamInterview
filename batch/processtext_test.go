package batch

import (
	"github.com/aalempijevic/communityteaminterview/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tokenizeWordsTests = []struct {
	text          string
	expectedWords []string
}{
	{"simple test", []string{"simple", "test"}},
	{"punctuation is removed!?!", []string{"punctuation", "is", "removed"}},
	{"Words are LOWER CASE", []string{"words", "are", "lower", "case"}},
	{"Contractions don't interfere", []string{"contractions", "dont", "interfere"}},
}

func TestTokenizeWords(t *testing.T) {
	for _, test := range tokenizeWordsTests {
		assert.Equal(t, test.expectedWords, tokenizeWords(test.text))
	}
}

var filterStopWordsTests = []struct {
	inputWords    []string
	expectedWords []string
}{
	{[]string{"simple", "test"}, []string{"simple", "test"}},
	{[]string{"the", "we", "are", "i"}, []string{}},
	{[]string{"some", "i", "the", "stop", "words"}, []string{"stop", "words"}},
}

func TestFilterStopWords(t *testing.T) {
	for _, test := range filterStopWordsTests {
		assert.Equal(t, test.expectedWords, filterStopWords(test.inputWords))
	}
}

var sentenceProcessingTests = []struct {
	sentences                  model.Sentences
	expectedWordSet            model.WordSet
	expectedTagWordFrequencies model.TagWordFrequencies
}{
	{
		model.Sentences{
			model.Sentence{TagIds: []int{1, 3}, Text: "sentence sentence poorly #@^ A. poorly; !punctuated !@#%sentence"},
			model.Sentence{TagIds: []int{1, 2}, Text: "Another sentence with Words"},
		},
		model.WordSet{
			"sentence":   true,
			"poorly":     true,
			"punctuated": true,
			"another":    true,
			"words":      true,
		},
		model.TagWordFrequencies{
			1: model.WordFrequencies{
				"sentence":   4,
				"poorly":     2,
				"punctuated": 1,
				"another":    1,
				"words":      1,
			},
			2: model.WordFrequencies{
				"sentence": 1,
				"another":  1,
				"words":    1,
			},
			3: model.WordFrequencies{
				"sentence":   3,
				"poorly":     2,
				"punctuated": 1,
			},
		},
	},
}

func TestProcessSentences(t *testing.T) {
	for _, test := range sentenceProcessingTests {
		wordSet, tagWordFrequencies := processSentences(test.sentences)

		assert.Equal(t, test.expectedWordSet, wordSet)
		assert.Equal(t, test.expectedTagWordFrequencies, tagWordFrequencies)
	}
}
