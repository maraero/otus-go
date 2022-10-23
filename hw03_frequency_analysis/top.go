package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const maxNumber = 10

func Top10(input string) []string {
	freqByWords := calcWords(input)
	wordFreqPairs := getWordFrequencyPairs(freqByWords)
	return getMostFrequentWordsByLimit(wordFreqPairs, maxNumber)
}

func calcWords(text string) map[string]int {
	words := strings.Fields(text)
	freqByWords := make(map[string]int)

	for _, word := range words {
		freqByWords[word]++
	}

	return freqByWords
}

type WordFreqPair struct {
	freq int
	word string
}

type FrequencyComparator []WordFreqPair

func (a FrequencyComparator) Len() int      { return len(a) }
func (a FrequencyComparator) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FrequencyComparator) Less(i, j int) bool {
	if a[i].freq == a[j].freq {
		return a[i].word > a[j].word
	}

	return a[i].freq < a[j].freq
}

func getWordFrequencyPairs(freqByWords map[string]int) []WordFreqPair {
	wordFreqPairs := []WordFreqPair{}
	for word, freq := range freqByWords {
		pair := WordFreqPair{freq: freq, word: word}
		wordFreqPairs = append(wordFreqPairs, pair)
	}
	sort.Sort(FrequencyComparator(wordFreqPairs))
	return wordFreqPairs
}

func getMostFrequentWordsByLimit(list []WordFreqPair, limit int) []string {
	result := []string{}
	for i := len(list) - 1; i >= 0 && limit > 0; i-- {
		result = append(result, list[i].word)
		limit--
	}
	return result
}
