package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Top10(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	reg := regexp.MustCompile(`^[[:punct:]]+|[[:punct:]]+$`)

	frequency := map[string]int{}

	for _, word := range words {
		if isPunctuationOnly(word) {
			if utf8.RuneCountInString(word) == 1 {
				continue
			}
		} else {
			word = reg.ReplaceAllString(word, "")
		}

		frequency[word]++
	}

	return sortTop(frequency)
}

func isPunctuationOnly(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

func sortTop(frequency map[string]int) []string {
	type keyValue struct {
		key   string
		value int
	}

	const defaultTop = 10

	var sorted []keyValue
	for key, value := range frequency {
		sorted = append(sorted, keyValue{key, value})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].value == sorted[j].value {
			return sorted[i].key < sorted[j].key
		}
		return sorted[i].value > sorted[j].value
	})

	end := defaultTop
	if len(sorted) < end {
		end = len(sorted)
	}

	sorted = sorted[:end]

	result := make([]string, 0, len(sorted))
	for _, s := range sorted {
		result = append(result, s.key)
	}

	return result
}
