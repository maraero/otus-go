package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

const maxNumber = 10

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

func Top10(input string) []string {
	freqByWords := calcWords(input)
	wordFreqPairs := make([]WordFreqPair, len(freqByWords))
	for word, freq := range freqByWords {
		pair := WordFreqPair{freq: freq, word: word}
		wordFreqPairs = append(wordFreqPairs, pair)
	}
	fmt.Println(wordFreqPairs)

	wordsByFreq := make(map[int][]string)

	for word, freq := range freqByWords {
		wordsByFreq[freq] = append(wordsByFreq[freq], word)
	}

	freqList := []int{}

	for freq := range wordsByFreq {
		freqList = append(freqList, freq)
	}

	sort.Slice(freqList, func(i, j int) bool { return freqList[i] > freqList[j] })
	result := []string{}

	for _, freq := range freqList {
		words := wordsByFreq[freq]
		sort.Slice(words, func(i, j int) bool { return words[i] < words[j] })

		for _, word := range words {
			result = append(result, word)

			if len(result) == maxNumber {
				return result
			}
		}
	}

	return result
}
