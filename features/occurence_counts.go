package features


// We want features that give the index of the nth occurrence of each letter in our alphabet in the
// input string. For example, for the string "edcba" we would have:
//    firstOccurrences = [4, 3, 2, 1, 0, ...]
// When the character doesn't occur in the string we say that the nth occurrence is \infty
// (analogous to stopping-times), which in our case would be 255. So to complete the above example
// we would have:
//    firstOccurrences = [4, 3, 2, 1, 0, 255, 255, ...]


import "encoding/gob"


//----------------------------------------------------------------------------------------------------
// Wrapper around OccurrenceCounts to provide FeatureSetConfig
var DefaultOccurrenceCounts OccurrenceCounts

type OccurrenceCounts struct {
    DirectionIsHead bool
    NumberOfOccurrences byte
}

func (o OccurrenceCounts) Size() int32 {
    return int32(ALPHABET_SIZE) * int32(o.NumberOfOccurrences)
}

func (o OccurrenceCounts) FromStringInPlace(input string, featureArray []byte) {
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
    allCharCounts := make([]byte, ALPHABET_SIZE)
    processChar := func(posInString int, ch byte) {
        charIndex := CHAR_MAP[ch]
        charCount := allCharCounts[charIndex]
        if charCount < o.NumberOfOccurrences {
            featureArray[(charCount * byte(ALPHABET_SIZE)) + byte(charIndex)] = byte(posInString)
            allCharCounts[charIndex] += 1
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
            ch := sNormalized[trimmedSMaxPos - i]
            processChar(i, ch)
        }
    }
}
//----------------------------------------------------------------------------------------------------


func init() {
    gob.Register(OccurrenceCounts{})
    DefaultOccurrenceCounts = OccurrenceCounts{true, 3}
}