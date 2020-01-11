package features

import "testing"

func TestGetFirstNumber(t *testing.T) {
	// given:
	strs := []string{"123 main st 789", "1st st 456 789", "abcd 234 main st 789", "main st", "main st 345"}
	expectedNums := []byte{123, 456 % 256, 234, 0, 345 % 256}
	// when/then:
	for i := 0; i < 4; i++ {
		n := GetFirstNumber(strs[i])
		if n != expectedNums[i] {
			t.Errorf("bad first-num %d for string %s, expected %d\n", n, strs[i], expectedNums[i])
		}
	}
}

func TestGetFirstNumberFeature(t *testing.T) {
	// given:
	fn := DefaultFirstNumber
	_, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{fn})
	strs := []string{"123 main st 789", "abcd 234 main st 789", "main st", "main st 345"}
	expectedNums := []byte{123, 234, 0, 345 % 256}
	// when:
	featureArray := fromStringArray(strs)
	// then:
	for i := 0; i < 4; i++ {
		if !testSliceIsSingleValue(featureArray[i], expectedNums[i]) {
			n := expectedNums[i]
			t.Errorf("bad feature-array for string %s, expected [%d, %d, %d], got %v\n", strs[i], n, n, n, featureArray[i])
		}
	}
}
