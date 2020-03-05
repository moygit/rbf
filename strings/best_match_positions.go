package strings

// Calculating edit distance between strings of length m and n is O(mn).
// In our case we also don't want distance against the full query string,
// we only want it against the "best" substring.

// The function GetBestMatchPositions below finds the "best" candidate
// substrings in time O(m+n). "Best" here just means that they have mostly
// the same letters (above some threshold) as the reference string.

import (
	"math"
	gostrings "strings"
)

const match_float_threshold = 0.67

var exclude_start_chars = make(map[byte]bool)

// Given a reference string, match its character counts against sliding windows of the query
// string. If the counts match above a certain threshold (0.67) then we'll (separately) do a
// Levenshtein match against those windows of the query string.
func GetBestMatchPositions(origReference, origQuery string) []int {
	reference := []byte(ConvertSpecialCharsToSpace(gostrings.ToLower(origReference)))
	query := []byte(ConvertSpecialCharsToSpace(gostrings.ToLower(origQuery)))

	intersectionCounts := getIntersections(reference, query)
	threshold := int(math.Round(match_float_threshold * float64(len(reference))))
	return getLocalMaximaAboveThreshold(intersectionCounts, threshold)
}

func init() {
	exclude_start_chars[' '] = true
}

// We'll slide the reference string along the query string, getting intersection counts at each position
// as we go. We only need to do the full calculation the first time; subsequent calculations only require
// us to drop the old first char and add the new last char.
func getIntersections(reference, query []byte) (intersectionTracker []int) {
	refCounts := getCounts(reference)
	lenRef := len(reference)

	// doing the 0th iteration of the loop outside b/c there's no dropChar/addChar-handling:
	startPos := 0
	endPos := min2i(len(query), lenRef)
	windowCounts := getCounts(query[:endPos])
	intersection := getIntersection(refCounts, windowCounts)
	lastIntersection := 0
	intersectionTracker = appendCheckedCount(intersectionTracker, intersection, lastIntersection, query, startPos)

	for startPos, endPos = startPos+1, endPos+1; endPos <= len(query); startPos, endPos = startPos+1, endPos+1 {
		// sliding the window, so drop one char and add a new one:
		dropChar := char_map[query[startPos-1]]
		addChar := char_map[query[endPos-1]]
		// remove dropped char from intersection:
		if windowCounts[dropChar] -= 1; windowCounts[dropChar] < refCounts[dropChar] {
			intersection -= 1
		}
		// add new char to intersection:
		if windowCounts[addChar] += 1; windowCounts[addChar] <= refCounts[addChar] {
			intersection += 1
		}
		intersectionTracker = appendCheckedCount(intersectionTracker, intersection, lastIntersection, query, startPos)
		lastIntersection = intersection
	}

	return
}

func getLocalMaximaAboveThreshold(arr []int, threshold int) []int {
	positions := make([]int, 0)
	last := len(arr) - 1
	for i := 0; i < len(arr); i++ {
		if (i == 0 || arr[i] >= arr[i-1]) && (i == last || arr[i] >= arr[i+1]) && (arr[i] >= threshold) {
			positions = append(positions, i)
		}
	}
	return positions
}

func getIntersection(refCounts, queryCounts []int) int {
	intersection := 0
	for i := 0; i < len(refCounts); i++ {
		intersection += min2i(refCounts[i], queryCounts[i])
	}
	return intersection
}

func appendCheckedCount(intersectionTracker []int, intersection, lastIntersection int, query []byte, startPos int) []int {
	if startPos < len(query) && !exclude_start_chars[query[startPos]] {
		return append(intersectionTracker, intersection)
	} else {
		return append(intersectionTracker, lastIntersection)
	}
}

func getCounts(s []byte) []int {
	counts := make([]int, alphabet_size)
	for _, ch := range s {
		counts[char_map[ch]] += 1
	}
	return counts
}

func min2i(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func min(x ...int) int {
	minVal := math.MaxInt64
	for _, i := range x {
		if i < minVal {
			minVal = i
		}
	}
	return minVal
}

func max(x ...int) int {
	maxVal := math.MinInt64
	for _, i := range x {
		if i > maxVal {
			maxVal = i
		}
	}
	return maxVal
}
