package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"unicode"
)

func main() {
	// Read text file
	text, err := ioutil.ReadFile("sample.txt")
	if err != nil {
		log.Fatalln("An error occurred while reading text file:", err.Error())
	}

	// Make text normalization
	normalizedText := normalizeText(text)

	fmt.Println("Top 10 words in text:")
	fmt.Println("---------------------")

	// Count and print top 10 words
	words := top10(normalizedText)
	for i, word := range words {
		fmt.Printf("%d: %s\n", i, word)
	}
}

func top10(text string) []string {
	// Split text by space
	strs := strings.Fields(text)

	// Create and fill map with words
	hist := make(map[string]int)
	for _, word := range strs {
		// Skip articles and unions
		if 3 > len(word) {
			continue
		}
		hist[word]++
	}

	type freqPair struct {
		word      string
		frequency int
	}
	frequencies := make([]freqPair, 0, 1024)

	// Sort words by frequency
	for key, value := range hist {
		frequencies = append(frequencies, freqPair{word: key, frequency: value})
	}
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].frequency > frequencies[j].frequency
	})

	// Get top 10 words
	words := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(frequencies); i++ {
		words = append(words, frequencies[i].word)
	}
	return words
}

func normalizeText(text []byte) (normalizedText string) {
	// Convert byte array to rune array
	runes := []rune(string(text))
	for _, r := range runes {
		switch {
		case unicode.IsLetter(r):
			normalizedText += string(unicode.ToLower(r))

		case unicode.IsSpace(r):
			normalizedText += string(r)
		}
	}
	return normalizedText
}
