package features

import (
	"encoding/binary"
	"io"
)

//----------------------------------------------------------------------------------------------------
// Provide FeatureSetConfig
var BigramsWithRepeats Bigrams
var BigramsNoRepeats Bigrams

type Bigrams struct {
	maxBigramCount uint8 // 255 if AllowRepeats else 1
	AllowRepeats   bool
}

func (b Bigrams) Size() int32 {
	return int32(alphabet_size * alphabet_size)
}

func (b Bigrams) FromStringInPlace(input string, featureArray []byte) {
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

const bigrams_type = int32(51)

func (b Bigrams) Serialize(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, bigrams_type)
	binary.Write(writer, binary.LittleEndian, int32(b.maxBigramCount))
}

func deserializeBigramsMap(confMap map[string]string) (config FeatureSetConfig, ok bool) {
	if allowRepeats, ok := confMap["allow_repeats"]; ok {
		maxBigramCount := 1
		if allowRepeats == "true" {
			maxBigramCount = 255
		}
		return Bigrams{maxBigramCount: byte(maxBigramCount), AllowRepeats: allowRepeats == "true"}, true
	}
	return nil, false
}

func deserialize_bigrams(reader io.Reader) FeatureSetConfig {
	var maxBigramCount int32
	binary.Read(reader, binary.LittleEndian, &maxBigramCount)
	return Bigrams{maxBigramCount: byte(maxBigramCount), AllowRepeats: maxBigramCount == 255}
}

//----------------------------------------------------------------------------------------------------

func init() {
	BigramsWithRepeats = Bigrams{maxBigramCount: 255, AllowRepeats: true}
	BigramsNoRepeats = Bigrams{maxBigramCount: 1, AllowRepeats: false}
}

func (b Bigrams) fromString(input string) []byte {
	featureArray := make([]byte, alphabet_size*alphabet_size)
	b.FromStringInPlace(input, featureArray)
	return featureArray
}
