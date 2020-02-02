package features

//----------------------------------------------------------------------------------------------------
// Provide featureSetConfig
type bigrams struct {
	maxBigramCount uint8 // 255 if we allow repeats else 1
}

func (b bigrams) Size() int32 {
	return int32(alphabet_size * alphabet_size)
}

func (b bigrams) FromStringInPlace(input string, featureArray []byte) {
	input = normalizeString(input)
	inputLen := len(input)

	for i := 0; i < inputLen-1; i++ {
		ch1 := char_map[input[i]]
		ch2 := char_map[input[i+1]]
		bigramIndex := (ch1 * alphabet_size) + ch2
		currentCount := featureArray[bigramIndex]
		if currentCount < b.maxBigramCount {
			featureArray[bigramIndex] = currentCount + 1
		}
	}
}

func deserializeBigramsMap(confMap map[string]string) (config featureSetConfig, ok bool) {
	if allowRepeats, ok := confMap["allow_repeats"]; ok {
		maxBigramCount := 1
		if allowRepeats == "true" {
			maxBigramCount = 255
		}
		return bigrams{maxBigramCount: byte(maxBigramCount)}, true
	}
	return nil, false
}

//----------------------------------------------------------------------------------------------------

func (b bigrams) fromString(input string) []byte {
	featureArray := make([]byte, alphabet_size*alphabet_size)
	b.FromStringInPlace(input, featureArray)
	return featureArray
}
