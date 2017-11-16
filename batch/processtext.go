package batch

import (
	"github.com/aalempijevic/communityteaminterview/model"
	"github.com/aalempijevic/communityteaminterview/repository"
	"log"
	"regexp"
	"strings"
)

//Punctuation is a regexp representing all non alphanumeric and whitespace characters
var punctuation = regexp.MustCompile("[^a-zA-Z0-9\\s]")

//TruncateAndProcessWords will truncate the tables for words and tag word frequencies and
// then populate them with counts for words in sentences scored higher than the threshold
func TruncateAndProcessWords(wordRepo repository.WordRepo, threshold float64) {

	wordRepo.Truncate()

	sentences, err := wordRepo.GetSentenceTags(threshold)
	if err != nil {
		log.Fatalf("Exception loading sentences from repo: %s", err)
	}

	log.Printf("There are %v sentences that meet the threshold %v \n", len(sentences), threshold)
	wordSet, tagWordFrequencies := processSentences(sentences)

	for k, v := range tagWordFrequencies {
		log.Printf("Tag id %v has %v unique words\n", k, len(v))
	}

	log.Println("Inserts are not batched so this might take a few minutes...")
	log.Println("Writing words to wordRepo")
	if err := wordRepo.StoreWords(wordSet); err != nil {
		log.Fatal(err)
	}

	log.Println("Writing tag word frequencies to wordRepo")
	if err := wordRepo.StoreWordFrequencies(tagWordFrequencies); err != nil {
		log.Fatal(err)
	}
	log.Println("Completed")
}

//ProcessSentences tokenizes sentences and counts their words to each tag they have applied
func processSentences(sentences model.Sentences) (model.WordSet, model.TagWordFrequencies) {
	tagWordFrequencies := model.TagWordFrequencies{}
	wordSet := model.WordSet{}

	for _, sentence := range sentences {
		words := filterStopWords(tokenizeWords(sentence.Text))

		wordSet.Append(words)

		for _, tagId := range sentence.TagIds {
			tagWordFrequencies.Append(tagId, words)
		}
	}

	return wordSet, tagWordFrequencies
}

//TokenizeWords returns slice of lower cased words in string without punctuation
func tokenizeWords(text string) []string {
	return strings.Fields(punctuation.ReplaceAllString(strings.ToLower(text), ""))
}

//FilterStopWords returns a slice of strings with configured stop words removed from input slice
func filterStopWords(words []string) []string {
	filteredWords := make([]string, 0)
	for _, word := range words {
		if !stopWords[word] {
			filteredWords = append(filteredWords, word)
		}
	}
	return filteredWords
}
