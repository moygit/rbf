package features

// We want features that give the count of each character in our alphabet in the
// input string. For example, for the string "aaabbcddd" we would have:
//    counts == [3, 2, 1, 3, 0, ...] ('a' occurs 3 times, 'b' 2 times, etc.)
// NOTE: We allow the user to specify the number of times they want this feature repeated
// (poor man's weighting).

import "strconv"

//----------------------------------------------------------------------------------------------------
// Provide featureSetConfig
type occurrenceCounts struct {
	Count byte
}

func (o occurrenceCounts) Size() int32 {
	return int32(alphabet_size) * int32(o.Count)
}

func (o occurrenceCounts) FromStringInPlace(input string, featureArray []byte) {
	sNormalized := []byte(normalizeString(input))
	for _, ch := range sNormalized {
		charIndex := char_map[ch]
		for i := 0; i < int(o.Count); i++ {
			currentCount := featureArray[i*alphabet_size+charIndex]
			if currentCount < 255 {
				featureArray[i*alphabet_size+charIndex] = currentCount + 1
			}
		}
	}
}

func deserializeOccurrenceCountsMap(confMap map[string]string) (config featureSetConfig, ok bool) {
	if countStr, ok := confMap["count"]; ok {
		if count, err := strconv.Atoi(countStr); err == nil {
			return occurrenceCounts{byte(count)}, true
		}
	}
	return nil, false
}
