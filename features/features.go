package features

// Interface for different types of features.
//
// Example usage:
//   featureSetConfigs := []features.FeatureSetConfig{ features.Followgrams{20}, features.FirstNumber{10}, features.DefaultLastNumber }
//   calculateFeatures, calculateFeaturesForArray := features.MakeFeatureCalculationFunctions(featureSetConfigs)
// And you can then use the `calculateFeatures` and `calculateFeaturesForArray` functions
// to calculate features for either a single string or an array of strings.
//   features := calculateFeatures("abcd")
//   // features is now an array that contains followgrams, first-number features, and last-number features for "abcd"
//   featuresArray := calculateFeaturesForArray([]string{"abcd", "efgh"})
//   // featuresArray[0] contains followgrams, first-number features, and last-number features for "abcd"
//   // featuresArray[1] contains followgrams, first-number features, and last-number features for "efgh"

import (
	"encoding/binary"
	"io"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"log"
)

type FeatureSetConfig interface {
	// can't just use an int here because we want to serialize this
	// and go serialization doesn't handle unsized ints well
	Size() int32

	// Given the input string s, put features for s into the given byte-slice.
	// Note: we do no position or size checking on the slice.
	FromStringInPlace(s string, features []byte)

	// write a type-identifier and then write data needed to reconstruct the value
	Serialize(writer io.Writer)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789 "

var alphabet_size int
var char_map map[byte]int
var non_alnum_pattern *regexp.Regexp

// TODO: remove
var CHAR_REVERSE_MAP map[int32]string

// Internally, we use a feature-set config to build a "featureSetRealized", which has the
// feature-generating function from the feature-set config and also has the start and end
// positions of this feature-set in a feature-array.
type featureSetRealized struct {
	start             int
	end               int
	fromStringInPlace func(s string, features []byte)
}

func init() {
	alphabet_size = len(alphabet)
	non_alnum_pattern = regexp.MustCompile("[^a-z0-9]+")

	char_map = make(map[byte]int)
	for i := 0; i < alphabet_size; i++ {
		char_map[alphabet[i]] = i
	}

	CHAR_REVERSE_MAP = make(map[int32]string)
	for i := 0; i < alphabet_size; i++ {
		for j := 0; j < alphabet_size; j++ {
			CHAR_REVERSE_MAP[int32((i*alphabet_size)+j)] = alphabet[i:i+1] + alphabet[j:j+1]
		}
	}
}

// For example usage please see package godoc above.
func MakeFeatureCalculationFunctions(selectedFeatureSetConfigs []FeatureSetConfig) (func(string) []byte, func([]string) [][]byte) {
	// Calculate feature-set sizes and positions in feature-array
	featureDefinitions := make([]featureSetRealized, len(selectedFeatureSetConfigs))
	var startPos int
	var totalNumFeatures int
	for i, featureSetConfig := range selectedFeatureSetConfigs {
		thisFeatureSetSize := int(featureSetConfig.Size())
		totalNumFeatures += thisFeatureSetSize
		featureDefinitions[i] = featureSetRealized{startPos, totalNumFeatures, featureSetConfig.FromStringInPlace}
		startPos += thisFeatureSetSize
	}

	// Given an input string and a byte slice to contain features, calculate the features
	// from each contained feature-set and put them in the appropriate place in the byte slice.
	fromStringInPlace := func(input string, features []byte) {
		for _, feature := range featureDefinitions {
			feature.fromStringInPlace(input, features[feature.start:feature.end])
		}
	}

	fromString := func(input string) []byte {
		features := make([]byte, totalNumFeatures)
		fromStringInPlace(input, features)
		return features
	}

	fromStringArray := func(inputArray []string) [][]byte {
		featuresArray2D := make([][]byte, len(inputArray))
		flattenedFeaturesArray := make([]byte, len(inputArray)*totalNumFeatures)
		for i, input := range inputArray {
			featuresArray2D[i] = flattenedFeaturesArray[(i * totalNumFeatures):((i + 1) * totalNumFeatures)]
			fromStringInPlace(input, featuresArray2D[i])
		}
		return featuresArray2D
	}

	return fromString, fromStringArray
}

func normalizeString(s string) string {
	return non_alnum_pattern.ReplaceAllLiteralString(strings.ToLower(s), " ")
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
		features[i] = deserialize(reader)
	}
	return features
}

func NewFeatureCalculationFunctionsFromConfig(confStr string) (func(string) []byte, func([]string) [][]byte) {
	configs := getConfigsFromYaml(confStr)
	return MakeFeatureCalculationFunctions(configs)
}

func getConfigsFromYaml(confStr string) (configs []FeatureSetConfig) {
	confMaps := make([]map[string]string, 0, 256)
	if err := yaml.Unmarshal([]byte(confStr), &confMaps); err != nil {
		log.Panicf("Error in feature-set yaml config: %v", err)
	}
	configs = make([]FeatureSetConfig, 0, 256)
	for _, confMap := range confMaps {
		confMap := mapToLowercase(confMap)
		configs = append(configs, deserializeMap(confMap))
	}
	return
}

func mapToLowercase(inMap map[string]string) (outMap map[string]string) {
	outMap = make(map[string]string, len(inMap))
	for key, val := range inMap {
		outMap[strings.ToLower(key)] = strings.ToLower(val)
	}
	return
}

func deserializeMap(confMap map[string]string) (config FeatureSetConfig) {
	var type_ string
	var ok bool
	if type_, ok = confMap["feature_type"]; !ok {
		log.Panicf("Feature config in yaml does not contain key 'feature_type': %v", confMap)
	}
	switch type_ {
	case "bigrams":
		config, ok = deserializeBigramsMap(confMap)
	case "followgrams":
		config, ok = deserializeFollowgramsMap(confMap)
	case "first_number":
		config, ok = deserializeFirstNumberMap(confMap)
	case "last_number":
		config, ok = deserializeLastNumberMap(confMap)
	case "occurrence_positions":
		config, ok = deserializeOccurrencePositionsMap(confMap)
	case "occurrence_counts":
		config, ok = deserializeOccurrenceCountsMap(confMap)
	default:
		log.Panicf("Received unknown feature-set type identifier " + type_)
	}
	if !ok {
		log.Panicf("Error deserializing feature-set config: %v", confMap)
	}
	return
}

func deserialize(reader io.Reader) FeatureSetConfig {
	var typeIdentifier int32
	binary.Read(reader, binary.LittleEndian, &typeIdentifier)
	switch typeIdentifier {
	case bigrams_type:
		return deserialize_bigrams(reader)
	case followgrams_type:
		return deserialize_followgrams(reader)
	case first_number_type:
		return deserialize_first_number(reader)
	case last_number_type:
		return deserialize_last_number(reader)
	case occurrence_positions_type:
		return deserialize_occurrence_positions(reader)
	case occurrence_counts_type:
		return deserialize_occurrence_counts(reader)
	default:
		panic("received unknown type identifier " + strconv.Itoa(int(typeIdentifier)))
	}
}

// Used for tests of most implementations
func testSliceIsSingleValue(slice []byte, val byte) bool {
	for _, sliceVal := range slice {
		if sliceVal != val {
			return false
		}
	}
	return true
}
