package features

// We want features that give the index of the nth occurrence of each letter in our alphabet in the
// input string. For example, for the string "edcba" we would have:
//    firstOccurrences == [4, 3, 2, 1, 0, ...] ('a' occurs in 4th pos, 'b' in 3rd pos, etc.)
// When the character doesn't occur in the string we say that the nth occurrence is \infty
// (analogous to stopping-times), which in our case would be 255. So to complete the above example
// we would have:
//    firstOccurrences == [4, 3, 2, 1, 0, 255, 255, ...]

import "strconv"

//----------------------------------------------------------------------------------------------------
// Provide featureSetConfig
type occurrencePositions struct {
	DirectionIsHead     bool
	NumberOfOccurrences byte
}

func (o occurrencePositions) Size() int32 {
	return int32(alphabet_size) * int32(o.NumberOfOccurrences)
}

func (o occurrencePositions) FromStringInPlace(input string, featureArray []byte) {
	// trim string to max length
	sNormalized := []byte(normalizeString(input))
	sLength := len(sNormalized)
	if sLength >= 256 {
		if o.DirectionIsHead {
			sNormalized = sNormalized[:256]
		} else {
			sNormalized = sNormalized[sLength-256:]
		}
	}

	// first set everything to infinity
	for i := 0; i < len(featureArray); i++ {
		featureArray[i] = 255
	}

	// function to update the feature-array if we've seen the ith byte fewer than NumberOfOccurrences times
	allCharPositions := make([]byte, alphabet_size)
	processChar := func(posInString int, ch byte) {
		charIndex := char_map[ch]
		charPosition := allCharPositions[charIndex]
		if charPosition < o.NumberOfOccurrences {
			featureArray[(charPosition*byte(alphabet_size))+byte(charIndex)] = byte(posInString)
			allCharPositions[charIndex] += 1
		}
	}

	// iterate either forwards or backwards
	if o.DirectionIsHead {
		for i, ch := range sNormalized {
			processChar(i, ch)
		}
	} else {
		trimmedSMaxPos := len(sNormalized) - 1
		// we're counting upwards but i will now be the position from the right:
		for i := 0; i <= trimmedSMaxPos; i++ {
			ch := sNormalized[trimmedSMaxPos-i]
			processChar(i, ch)
		}
	}
}

func deserializeOccurrencePositionsMap(confMap map[string]string) (config featureSetConfig, ok bool) {
	var directionIsHead bool
	if directionIsHeadStr, ok := confMap["direction_is_head"]; !ok {
		return nil, false
	} else {
		directionIsHead = (directionIsHeadStr == "true")
	}

	if numOccurrencesStr, ok := confMap["num_occurrences"]; !ok {
		return nil, false
	} else if numOccurrences, err := strconv.Atoi(numOccurrencesStr); err != nil {
		return nil, false
	} else {
		return occurrencePositions{directionIsHead, byte(numOccurrences)}, true
	}
}
