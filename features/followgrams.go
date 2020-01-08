package features


import (
    "encoding/binary"
    "io"
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

    for i := 0; i < sNormalizedLen - 1; i++ {
        ch1 := CHAR_MAP[sNormalized[i]]

        // get window right edge, making sure we don't fall off the end of the string
        followgramWindowEnd := i + f.WindowSize + 1
        if followgramWindowEnd > sNormalizedLen {
            followgramWindowEnd = sNormalizedLen
        }
        for j := i + 1; j < followgramWindowEnd; j++ {
            // get the index into the followgram array and increment the count,
            // making sure we don't overflow the byte
            ch2 := CHAR_MAP[sNormalized[j]]
            followgramIndex := (ch1 * ALPHABET_SIZE) + ch2
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

func deserialize_followgrams(reader io.Reader) FeatureSetConfig {
    var windowSize int32
    binary.Read(reader, binary.LittleEndian, &windowSize)
    return Followgrams{int(windowSize)}
}
//----------------------------------------------------------------------------------------------------


func init() {
    DefaultFollowgrams = Followgrams{followgram_default_window_size}
    num_followgrams = ALPHABET_SIZE * ALPHABET_SIZE
}


func (f Followgrams) FromString(input string) []byte {
    featureArray := make([]byte, num_followgrams)
    f.FromStringInPlace(input, featureArray)
    return featureArray
}
