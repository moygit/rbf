package features

import "testing"

func TestGetFirstNumber(t *testing.T) {
	// given:
	strs := []string{"123 main st 789", "1st st 456 789", "abcd 234 main st 789", "main st", "main st 345"}
	expectedNums := []byte{123, 456 % 256, 234, 0, 345 % 256}
	// when/then:
	for i, str := range strs {
		n := GetFirstNumber(str)
		if n != expectedNums[i] {
			t.Errorf("bad first-num %d for string %s, expected %d\n", n, str, expectedNums[i])
		}
	}
}

func TestGetFirstNumberAsString(t *testing.T) {
	// given:
	strs := []string{"000123 main st 789", "1st st 456 789", "abcd 234 main st 789", "main st", "main st 345"}
	expectedNums := []string{"123", "456", "234", "", "345"}
	// when/then:
	for i, str := range strs {
		n := GetFirstNumberAsString(str)
		if n != expectedNums[i] {
			t.Errorf("bad first-num %s for string %s, expected %s\n", n, str, expectedNums[i])
		}
	}
}

func TestGetFirstNumberFeature(t *testing.T) {
	// given:
	config := "- feature_type: first_number\n  count: 20"
	_, fromStringArray := CreateFeatureCalcFuncs(config)
	strs := []string{"123 main st 789", "1st st 456 789", "abcd 234 main st 789", "main st", "main st 345"}
	expectedNums := []byte{123, 456 % 256, 234, 0, 345 % 256}
	// when:
	featureArray := fromStringArray(strs)
	// then:
	for i, str := range strs {
		if !testSliceIsSingleValue(featureArray[i], expectedNums[i]) {
			n := expectedNums[i]
			t.Errorf("bad feature-array for string %s, expected [%d, %d, %d], got %v\n", str, n, n, n, featureArray[i])
		}
	}
}
