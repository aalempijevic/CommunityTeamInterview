package model

type WordSet map[string]bool

//Append will add word counts for words into the word map receiver
func (wordSet WordSet) Append(words []string) {
	for _, word := range words {
		if _, ok := wordSet[word]; !ok {
			wordSet[word] = true
		}
	}
}
