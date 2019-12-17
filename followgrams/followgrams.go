package followgrams


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
func FromString(s string) []byte {
    return getFollowgramsWithWindowSize(s, FOLLOWGRAM_DEFAULT_WINDOW_SIZE)
}


// Convert an array of ASCII strings to an array of followgram arrays.
func FromStringArray(sArray []string) [][]byte {
    followgramsArray := make([][]byte, len(sArray))
    fullFollowgramsArray := make([]byte, len(sArray) * NUM_FOLLOWGRAMS)
    for i, s := range sArray {
        followgramsArray[i] = fullFollowgramsArray[(i * NUM_FOLLOWGRAMS):((i + 1) * NUM_FOLLOWGRAMS)]
        getFollowgramsInPlace(s, followgramsArray[i])
    }
    return followgramsArray
}


func getFollowgramsWithWindowSize(s string, windowSize int) []byte {
    followgramsArray := make([]byte, NUM_FOLLOWGRAMS)
    getFollowgramsInPlaceWithWindowSize(s, windowSize, followgramsArray)
    return followgramsArray
}


func getFollowgramsInPlace(s string, followgramsArray []byte) {
    getFollowgramsInPlaceWithWindowSize(s, FOLLOWGRAM_DEFAULT_WINDOW_SIZE, followgramsArray)
}


func getFollowgramsInPlaceWithWindowSize(s string, windowSize int, followgramsArray []byte) {
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
