package features


// Example use:
// featureSetConfig := []FeatureSetConfig{ features.DefaultFollowgrams, features.DefaultOccurrenceCounts }


type FeatureSetConfig interface {
    Size() int32    // can't just use an int here because we want to serialize this
                    // and go serialization doesn't handle int's (without sizes) well
    FromStringInPlace(s string, features []byte)
}


type featureSetRealized struct {
    start int
    end int
    fromStringInPlace func(s string, features []byte)
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
