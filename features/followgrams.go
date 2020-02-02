package features

// Followgrams [same as k-skip n-grams]:
// - We know that an n-gram of a string is any n-length substring of that string.
// - By analogy, we define an "n-followgram" to be a pair "ab" such that "b" follows "a" in the parent
//   string within a window of size n+1 (e.g. the string "abcd" contains the 2-followgrams
//   ab, ac, (but not ad because that distance is 3), bc, bd, and cd).
// - We'll refer to the "infinity-followgrams" as just followgrams.
// - So any string, of any length, over the alphabet [a-z0-9 ] can have 1369 (37 * 37) different
//   followgrams: aa, ab, ..., az, a0, ..., a9, "a ", ba, bb, ..., "b ", ..., " a", " b", ..., "  ".
// - So for a string of length 256, the maximum count in the followgrams array is
//   255 + 254 + 253 + ...  + 1 = 255 * 254 / 2, roughly 2**16 - 1.
//   More generally, for a string of length 2**n, the max count in the followgrams array is roughly
//   2**(2n) - 1.
// - Notice that sum(followgrams(string)) must be (n)(n-1)/2, where n = len(string).
// - Given that, it's a fairly easy inductive proof that if s1 != s2 then
//   followgram(s1) != followgram(s2). (n-grams don't have this uniqueness property.)
//   Reconstruct the original string from the followgram array as follows: Sum the array grouping by
//   first letter; the letter with the highest sum must be the first letter of the string. Subtract 1
//   from each entry with that first letter, and now you have the followgram array of the rest of the
//   original string (i.e. original - first char).
// - In practice, I've found that this uniqueness property actually causes "infinity-followgrams" to
//   yield poor results in our use-case. So we'll use smaller followgram windows (say 5-followgrams) as
//   a proxy for using n>2-grams (say 6-grams).

import (
	"encoding/binary"
	"io"
	"strconv"
)

const followgram_default_window_size = 5
const max_followgram_count = 255

var num_followgrams int

//----------------------------------------------------------------------------------------------------
// Provide FeatureSetConfig
var DefaultFollowgrams Followgrams

type Followgrams struct {
	WindowSize int
}

func (f Followgrams) Size() int32 {
	return int32(num_followgrams)
}

func (f Followgrams) FromStringInPlace(input string, featureArray []byte) {
	sNormalized := normalizeString(input)
	sNormalizedLen := len(sNormalized)

	for i := 0; i < sNormalizedLen-1; i++ {
		ch1 := char_map[sNormalized[i]]

		// get window right edge, making sure we don't fall off the end of the string
		followgramWindowEnd := i + f.WindowSize + 1
		if followgramWindowEnd > sNormalizedLen {
			followgramWindowEnd = sNormalizedLen
		}
		for j := i + 1; j < followgramWindowEnd; j++ {
			// get the index into the followgram array and increment the count,
			// making sure we don't overflow the byte
			ch2 := char_map[sNormalized[j]]
			followgramIndex := (ch1 * alphabet_size) + ch2
			currentCount := featureArray[followgramIndex]
			if currentCount < max_followgram_count {
				featureArray[followgramIndex] = currentCount + 1
			}
		}
	}
}

const followgrams_type = int32(1)

func (f Followgrams) Serialize(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, followgrams_type)
	binary.Write(writer, binary.LittleEndian, int32(f.WindowSize))
}

func deserializeFollowgramsMap(confMap map[string]string) (config FeatureSetConfig, ok bool) {
	if windowSizeStr, ok := confMap["window_size"]; ok {
		if windowSize, err := strconv.Atoi(windowSizeStr); err == nil {
			return Followgrams{int(windowSize)}, true
		}
	}
	return nil, false
}

func deserialize_followgrams(reader io.Reader) FeatureSetConfig {
	var windowSize int32
	binary.Read(reader, binary.LittleEndian, &windowSize)
	return Followgrams{int(windowSize)}
}

//----------------------------------------------------------------------------------------------------

func init() {
	DefaultFollowgrams = Followgrams{followgram_default_window_size}
	num_followgrams = alphabet_size * alphabet_size
}

func (f Followgrams) fromString(input string) []byte {
	featureArray := make([]byte, num_followgrams)
	f.FromStringInPlace(input, featureArray)
	return featureArray
}
