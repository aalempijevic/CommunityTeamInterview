package model

//WordFrequencies is a map of word to frequency
type WordFrequencies map[string]int

//Append will add word counts for words into the word map receiver
func (wordMap WordFrequencies) Append(words []string) {
	for _, word := range words {
		if wordCount, ok := wordMap[word]; ok {
			wordMap[word] = wordCount + 1
		} else {
			wordMap[word] = 1
		}
	}
}
