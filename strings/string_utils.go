package strings


import (
    // "encoding/binary"
    // "fmt"
    // "io/ioutil"
    // "os"
    // "math/rand"
    "regexp"
    // "sort"
    "strings"

    // expand "github.com/openvenues/gopostal/expand"
    // parser "github.com/openvenues/gopostal/parser"

    // for logging only:
    // "log"
    // "os"
)


const ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789 "

const FOLLOWGRAM_DEFAULT_WINDOW_SIZE = 5
const MAX_FOLLOWGRAM_COUNT = 255

var ALPHABET_SIZE int
var NUM_FOLLOWGRAMS int

var CHAR_MAP map[byte]int
var NON_ALNUM_PATTERN *regexp.Regexp

// TODO: remove
var CHAR_REVERSE_MAP map[int]string


func init() {
    ALPHABET_SIZE = len(ALPHABET)
    NUM_FOLLOWGRAMS = ALPHABET_SIZE * ALPHABET_SIZE

    NON_ALNUM_PATTERN = regexp.MustCompile("[^a-z0-9]+")

    CHAR_MAP = make(map[byte]int)
    for i := 0; i < ALPHABET_SIZE; i++ {
        CHAR_MAP[ALPHABET[i]] = i
    }

    CHAR_REVERSE_MAP = make(map[int]string)
    for i := 0; i < ALPHABET_SIZE; i++ {
        for j := 0; j < ALPHABET_SIZE; j++ {
            CHAR_REVERSE_MAP[(i * ALPHABET_SIZE) + j] = ALPHABET[i:i+1] + ALPHABET[j:j+1]
        }
    }
}


// Normalize a string and get a followgrams array for the normalized version.
// Normalizing consists of:
// 1. lowercasing
// 2. replacing any sequence of characters not in a specified alphabet ([a-z0-9]) with a single whitespace
func GetFollowgrams(s string) []byte {
    return GetFollowgramsWithWindowSize(s, FOLLOWGRAM_DEFAULT_WINDOW_SIZE)
}

func GetFollowgramsWithWindowSize(s string, windowSize int) []byte {
    followgramsArray := make([]byte, NUM_FOLLOWGRAMS)
    GetFollowgramsInPlaceWithWindowSize(s, windowSize, followgramsArray)
    return followgramsArray
}


func GetFollowgramsInPlace(s string, followgramsArray []byte) {
    GetFollowgramsInPlaceWithWindowSize(s, FOLLOWGRAM_DEFAULT_WINDOW_SIZE, followgramsArray)
}


func GetFollowgramsInPlaceWithWindowSize(s string, windowSize int, followgramsArray []byte) {
    sNormalized := NON_ALNUM_PATTERN.ReplaceAllLiteralString(strings.ToLower(s), " ")
    sNormalizedLen := len(sNormalized)

    for i := 0; i < sNormalizedLen - 1; i++ {
        ch1 := CHAR_MAP[sNormalized[i]]

        // get window right edge, making sure we don't fall off the end of the string
        followgramWindowEnd := i + windowSize + 1
        if followgramWindowEnd > sNormalizedLen {
            followgramWindowEnd = sNormalizedLen
        }
        for j := i + 1; j < followgramWindowEnd; j++ {
            // get the index into the followgram array and increment the count,
            // making sure we don't overflow the byte
            ch2 := CHAR_MAP[sNormalized[j]]
            followgramIndex := (ch1 * ALPHABET_SIZE) + ch2
            currentCount := followgramsArray[followgramIndex]
            if currentCount < MAX_FOLLOWGRAM_COUNT {
                followgramsArray[followgramIndex] = currentCount + 1
            }
        }
    }
}


// Convert an array of ASCII strings to an array of followgram arrays.
func GetFollowgramsForArray(sArray []string) [][]byte {
    followgramsArray := make([][]byte, len(sArray))
    fullFollowgramsArray := make([]byte, len(sArray) * NUM_FOLLOWGRAMS)
    for i, s := range sArray {
        followgramsArray[i] = fullFollowgramsArray[(i * NUM_FOLLOWGRAMS):((i + 1) * NUM_FOLLOWGRAMS)]
        GetFollowgramsInPlace(s, followgramsArray[i])
    }
    return followgramsArray
}


func Jaccard2gramSimilarity(s1 string, s2 string) float64 {
    // Note: these are ASCIIized strings, so we can iterate over bytes instead of runes.

    allNGrams := make(map[int]bool)
    getNGramDict := func (s string) map[int]int {
        nGramDict := make(map[int]int)
        for i := 0; i < len(s) - 1; i++ {
            key := (CHAR_MAP[s[i]] * len(CHAR_MAP)) + CHAR_MAP[s[i+1]]  // treat each Ngram as a number
            nGramDict[key] += 1
            allNGrams[key] = true
        }
        return nGramDict
    }
    
    ngrams1 := getNGramDict(s1)
    ngrams2 := getNGramDict(s2)

    var intersection, union int
    for key, _ := range allNGrams {
        count1 := ngrams1[key]
        count2 := ngrams2[key]
        if count1 > count2 {
            intersection += count2
            union += count1
        } else {
            intersection += count1
            union += count2
        }
    }

    return float64(intersection) / float64(union)
}


func JaccardFollowgramSimilarity(s1 string, s2 string) float64 {
    // Note: these are ASCIIized strings, so we can iterate over bytes instead of runes.

    allFollowGrams := make(map[int]bool)
    getFollowGramDict := func (s string) map[int]int {
        followGramDict := make(map[int]int)
        for i := 0; i < len(s) - 1; i++ {
            for j := i + 1; j < len(s); j++ {
                key := (CHAR_MAP[s[i]] * len(CHAR_MAP)) + CHAR_MAP[s[j]]  // treat each Ngram as a number
                followGramDict[key] += 1
                allFollowGrams[key] = true
            }
        }
        return followGramDict
    }
    
    followgrams1 := getFollowGramDict(s1)
    followgrams2 := getFollowGramDict(s2)

    var intersection, union int
    for key, _ := range allFollowGrams {
        count1 := followgrams1[key]
        count2 := followgrams2[key]
        if count1 > count2 {
            intersection += count2
            union += count1
        } else {
            intersection += count1
            union += count2
        }
    }

    return float64(intersection) / float64(union)
}
