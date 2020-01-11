package strings

import (
	"regexp"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789 "

var char_map map[byte]int
var non_alnum_pattern *regexp.Regexp

func init() {
	char_map = make(map[byte]int)
	for i := 0; i < len(alphabet); i++ {
		char_map[alphabet[i]] = i
	}

	non_alnum_pattern = regexp.MustCompile("[^a-z0-9]+")
}

func LowercaseAndRemoveSpecialChars(s string) string {
	return non_alnum_pattern.ReplaceAllLiteralString(strings.ToLower(s), "")
}

func LowercaseAndSpaceSpecialChars(s string) string {
	return non_alnum_pattern.ReplaceAllLiteralString(strings.ToLower(s), " ")
}

// Experimental and inefficient, not in use right now and will likely be removed.
func OneDirectionalJaccard2GramSimilarity(reference string, eval string) float64 {
	// Note: these are ASCIIized strings, so we can iterate over bytes instead of runes.

	allNGrams := make(map[int]bool)
	getNGramDict := func(s string) map[int]int {
		nGramDict := make(map[int]int)
		for i := 0; i < len(s)-1; i++ {
			key := (char_map[s[i]] * len(char_map)) + char_map[s[i+1]] // treat each Ngram as a number
			nGramDict[key] += 1
			allNGrams[key] = true
		}
		return nGramDict
	}

	ngramReference := getNGramDict(reference)
	ngramEval := getNGramDict(eval)

	var intersection, union int
	totalRef := 0
	for key, _ := range allNGrams {
		countRef := ngramReference[key]
		totalRef += countRef
		countEval := ngramEval[key]
		if countRef > countEval {
			intersection += countEval
			union += countRef
		} else {
			intersection += countRef
			union += countEval
		}
	}

	return float64(intersection) / float64(totalRef)
}

// Experimental and inefficient, not in use right now and will likely be removed.
func Jaccard2gramSimilarity(s1 string, s2 string) float64 {
	// Note: these are ASCIIized strings, so we can iterate over bytes instead of runes.

	allNGrams := make(map[int]bool)
	getNGramDict := func(s string) map[int]int {
		nGramDict := make(map[int]int)
		for i := 0; i < len(s)-1; i++ {
			key := (char_map[s[i]] * len(char_map)) + char_map[s[i+1]] // treat each Ngram as a number
			nGramDict[key] += 1
			allNGrams[key] = true
		}
		return nGramDict
	}

	ngrams1 := getNGramDict(s1)
	ngrams2 := getNGramDict(s2)

	var intersection, union int
	for key, _ := range allNGrams {
		count1 := ngrams1[key]
		count2 := ngrams2[key]
		if count1 > count2 {
			intersection += count2
			union += count1
		} else {
			intersection += count1
			union += count2
		}
	}

	return float64(intersection) / float64(union)
}

// Experimental and inefficient, not in use right now and will likely be removed.
func JaccardFollowgramSimilarity(s1 string, s2 string) float64 {
	// Note: these are ASCIIized strings, so we can iterate over bytes instead of runes.

	allFollowGrams := make(map[int]bool)
	getFollowGramDict := func(s string) map[int]int {
		followGramDict := make(map[int]int)
		for i := 0; i < len(s)-1; i++ {
			for j := i + 1; j < len(s); j++ {
				key := (char_map[s[i]] * len(char_map)) + char_map[s[j]] // treat each Ngram as a number
				followGramDict[key] += 1
				allFollowGrams[key] = true
			}
		}
		return followGramDict
	}

	followgrams1 := getFollowGramDict(s1)
	followgrams2 := getFollowGramDict(s2)

	var intersection, union int
	for key, _ := range allFollowGrams {
		count1 := followgrams1[key]
		count2 := followgrams2[key]
		if count1 > count2 {
			intersection += count2
			union += count1
		} else {
			intersection += count1
			union += count2
		}
	}

	return float64(intersection) / float64(union)
}
