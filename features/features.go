package features


import (
    "encoding/binary"
    "io"
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
    Serialize(writer io.Writer)     // write a type-identifier and then write data needed to reconstruct the value
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


func testSliceIsSingleValue(slice []byte, val byte) bool {
    for _, sliceVal := range slice {
        if sliceVal != val {
            return false
        }
    }
    return true
}


func SerializeArray(features []FeatureSetConfig, writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, int32(len(features)))
    for _, feature := range features {
        feature.Serialize(writer)
    }
}

func DeserializeArray(reader io.Reader) []FeatureSetConfig {
    var length int32
    binary.Read(reader, binary.LittleEndian, &length)
    features := make([]FeatureSetConfig, length)
    for i := int32(0); i < length; i++ {
        features[i] = Deserialize(reader)
    }
    return features
}

func Deserialize(reader io.Reader) FeatureSetConfig {
    var typeIdentifier int32
    binary.Read(reader, binary.LittleEndian, &typeIdentifier)
    switch typeIdentifier {
    case followgrams_type:
        return deserialize_followgrams(reader)
    case first_number_type:
        return deserialize_first_number(reader)
    case last_number_type:
        return deserialize_last_number(reader)
    case occurrence_counts_type:
        return deserialize_occurrence_counts(reader)
    default:
        // TODO: error!
        return nil
    }
}
