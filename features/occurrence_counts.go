package features

// We want features that give the count of each character in our alphabet in the
// input string. For example, for the string "aaabbcddd" we would have:
//    counts == [3, 2, 1, 3, 0, ...] ('a' occurs 3 times, 'b' 2 times, etc.)
// NOTE: We allow the user to specify the number of times they want this feature repeated
// (poor man's weighting).

import (
	"encoding/binary"
	"io"
)

//----------------------------------------------------------------------------------------------------
// Provide FeatureSetConfig
var DefaultOccurrenceCounts OccurrenceCounts

type OccurrenceCounts struct {
	Count byte
}

func (o OccurrenceCounts) Size() int32 {
	return int32(alphabet_size) * int32(o.Count)
}

func (o OccurrenceCounts) fromStringInPlace(input string, featureArray []byte) {
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

const occurrence_counts_type = int32(31)

func (oc OccurrenceCounts) Serialize(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, occurrence_counts_type)
	binary.Write(writer, binary.LittleEndian, int32(oc.Count))
}

func deserialize_occurrence_counts(reader io.Reader) FeatureSetConfig {
	var count int32
	binary.Read(reader, binary.LittleEndian, &count)
	return OccurrenceCounts{byte(count)}
}

//----------------------------------------------------------------------------------------------------

func init() {
	DefaultOccurrenceCounts = OccurrenceCounts{2}
}
