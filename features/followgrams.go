package features


import (
    // "encoding/binary"
    // "fmt"
    // "io/ioutil"
    // "os"
    // "math/rand"
    "encoding/gob"
    // "sort"

    // expand "github.com/openvenues/gopostal/expand"
    // parser "github.com/openvenues/gopostal/parser"

    // for logging only:
    // "log"
    // "os"
)


const followgram_default_window_size = 5
const max_followgram_count = 255
var num_followgrams int


//----------------------------------------------------------------------------------------------------
// Wrapper around Followgrams to provide FeatureSetConfig
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
//----------------------------------------------------------------------------------------------------


func init() {
    gob.Register(Followgrams{})
    DefaultFollowgrams = Followgrams{followgram_default_window_size}
    num_followgrams = ALPHABET_SIZE * ALPHABET_SIZE
}


func (f Followgrams) FromString(input string) []byte {
    featureArray := make([]byte, num_followgrams)
    f.FromStringInPlace(input, featureArray)
    return featureArray
}
