package model

//WordFrequencies is a map of word to frequency
type WordFrequencies map[string]int

type TagWordFrequencies map[int]WordFrequencies

//Append will add word counts for words into the WordFrequencies map
func (wordFrequencies WordFrequencies) Append(words []string) {
	for _, word := range words {
		if wordCount, ok := wordFrequencies[word]; ok {
			wordFrequencies[word] = wordCount + 1
		} else {
			wordFrequencies[word] = 1
		}
	}
}

//Append will add word counts for words into the TagWordFrequencies map
func (tagWordFrequencies TagWordFrequencies) Append(tagId int, words []string) {
	if _, ok := tagWordFrequencies[tagId]; ok {
		tagWordFrequencies[tagId].Append(words)
	} else {
		tagWordFrequencies[tagId] = WordFrequencies{}
		tagWordFrequencies[tagId].Append(words)
	}
}
