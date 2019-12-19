package features


import (
    "regexp"
    "strings"
)


// Example use:
// featureSetConfig := []FeatureSetConfig{ features.DefaultFollowgrams, features.DefaultOccurrenceCounts }
//
// Internally, we use a feature-set config to build a "featureSetRealized", which has the
// feature-generating function from the feature-set config and also has the start and end
// positions of this feature-set in a feature-array.


type FeatureSetConfig interface {
    Size() int32    // can't just use an int here because we want to serialize this
                    // and go serialization doesn't handle int's (without sizes) well
    FromStringInPlace(s string, features []byte)
}


const ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789 "
var ALPHABET_SIZE int
var CHAR_MAP map[byte]int
var NON_ALNUM_PATTERN *regexp.Regexp

// TODO: remove
var CHAR_REVERSE_MAP map[int32]string


type featureSetRealized struct {
    start int
    end int
    fromStringInPlace func(s string, features []byte)
}


func init() {
    ALPHABET_SIZE = len(ALPHABET)
    NON_ALNUM_PATTERN = regexp.MustCompile("[^a-z0-9]+")

    CHAR_MAP = make(map[byte]int)
    for i := 0; i < ALPHABET_SIZE; i++ {
        CHAR_MAP[ALPHABET[i]] = i
    }

    CHAR_REVERSE_MAP = make(map[int32]string)
    for i := 0; i < ALPHABET_SIZE; i++ {
        for j := 0; j < ALPHABET_SIZE; j++ {
            CHAR_REVERSE_MAP[int32((i * ALPHABET_SIZE) + j)] = ALPHABET[i:i+1] + ALPHABET[j:j+1]
        }
    }
}


func MakeFeatureCalculationFunctions(selectedFeatureSetConfigs []FeatureSetConfig) (func (string) []byte, func([]string) [][]byte) {
    featureDefinitions := make([]featureSetRealized, len(selectedFeatureSetConfigs))
    var startPos int
    var totalNumFeatures int
    for i, featureSetConfig := range selectedFeatureSetConfigs {
        thisFeatureSetSize := int(featureSetConfig.Size())
        totalNumFeatures += thisFeatureSetSize
        featureDefinitions[i] = featureSetRealized{startPos, totalNumFeatures, featureSetConfig.FromStringInPlace}
        startPos += thisFeatureSetSize
    }

    fromStringInPlace := func(input string, features []byte) {
        for _, feature := range featureDefinitions {
            feature.fromStringInPlace(input, features[feature.start:feature.end])
        }
    }

    fromString := func (input string) []byte {
        features := make([]byte, totalNumFeatures)
        fromStringInPlace(input, features)
        return features
    }

    fromStringArray := func (inputArray []string) [][]byte {
        featuresArray2D := make([][]byte, len(inputArray))
        flattenedFeaturesArray := make([]byte, len(inputArray) * totalNumFeatures)
        for i, input := range inputArray {
            featuresArray2D[i] = flattenedFeaturesArray[(i * totalNumFeatures):((i + 1) * totalNumFeatures)]
            fromStringInPlace(input, featuresArray2D[i])
        }
        return featuresArray2D
    }

    return fromString, fromStringArray
}


func normalizeString(s string) string {
    return NON_ALNUM_PATTERN.ReplaceAllLiteralString(strings.ToLower(s), " ")
}
