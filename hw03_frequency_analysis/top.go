package hw03frequencyanalysis

import (
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

func Top10(input string) []string {
	freqByWords := calcWords(input)

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
