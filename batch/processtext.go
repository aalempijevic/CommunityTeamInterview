package batch

import (
	"strings"
	"regexp"
	"github.com/aalempijevic/communityteaminterview/repository"
	"github.com/aalempijevic/communityteaminterview/model"
	"github.com/reiver/go-porterstemmer"
	"log"
)

//Punctuation is a regexp representing all non alphanumeric and whitespace characters
var punctuation  = regexp.MustCompile("[^a-zA-Z0-9\\s]")

//TruncateAndProcessWords will truncate the tables for words and tag word frequencies and
// then populate them with counts for words in sentences scored higher than the threshold
func TruncateAndProcessWords(repo repository.WordRepo, threshold float32, useStemming bool) {
	tagWordFrequencies := make(map[int]model.WordFrequencies)
	wordSet := make(map[string]bool)
	//repo.TruncateWords()
	//repo.TruncateWordCounts()

	sentences := repo.GetSentenceTags(threshold)

	log.Printf("There are %v sentences that meet the threshold %v \n", len(sentences), threshold)
	for _, sentence := range sentences {
		words := filterStopWords(words(sentence.Text))
		if useStemming {
			words = stemWords(words)
		}
		for _, word := range words {
			if _, ok := wordSet[word]; !ok {
				wordSet[word] = true
			}
		}

		for _, tagId := range sentence.TagIds {
			if _, ok := tagWordFrequencies[tagId]; ok {
				tagWordFrequencies[tagId].Append(words)
			} else {
				tagWordFrequencies[tagId] = model.WordFrequencies{}
				tagWordFrequencies[tagId].Append(words)
			}

		}
	}

	for k, v := range tagWordFrequencies {
		log.Printf("Tag id %v has %v unique words\n", k, len(v) )
	}

	log.Println("Writing words to repo")
	if err := repo.StoreWords(wordSet); err != nil {
		log.Fatal(err)
	}

	log.Println("Writing tag word frequencies to repo")
	if err := repo.StoreWordFrequencies(tagWordFrequencies); err != nil {
		log.Fatal(err)
	}
}

func words(text string) []string {
	return strings.Fields(punctuation.ReplaceAllString(strings.ToLower(text), ""))
}

func filterStopWords(words []string) []string {
	filteredWords := make([]string, 0)
	for _, word := range words {
		if !stopWords[word] {
			filteredWords = append(filteredWords, word)
		}
	}
	return filteredWords
}

func stemWords(words []string) []string {
	stemmedWords := make([]string, len(words))
	for _, word := range words {
		stemmedWords = append(stemmedWords, porterstemmer.StemString(word))
	}
	return stemmedWords
}
