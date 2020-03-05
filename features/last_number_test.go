package features

import "testing"

func TestGetLastNumber(t *testing.T) {
	// given:
	strs := []string{"123 main st 789--", "123 456 1st st", "abcd 234 main st 678", "main st", "123 main st"}
	expectedNums := []byte{789 % 256, 456 % 256, 678 % 256, 0, 123}
	// when/then:
	for i := 0; i < 4; i++ {
		n := GetLastNumber(strs[i])
		if n != expectedNums[i] {
			t.Errorf("bad last-num %d for string %s, expected %d\n", n, strs[i], expectedNums[i])
		}
	}
}

func TestGetLastNumberAsString(t *testing.T) {
	// given:
	strs := []string{"123 main st 789--", "123 456 1st st", "abcd 234 main st 678", "main st", "123 main st"}
	expectedNums := []string{"789", "456", "678", "", "123"}
	// when/then:
	for i := 0; i < 4; i++ {
		n := GetLastNumberAsString(strs[i])
		if n != expectedNums[i] {
			t.Errorf("bad last-num %s for string %s, expected %s\n", n, strs[i], expectedNums[i])
		}
	}
}

func TestGetLastNumberFeature(t *testing.T) {
	// given:
	config := "- feature_type: last_number\n  count: 10"
	_, fromStringArray := CreateFeatureCalcFuncs(config)
	strs := []string{"123 main st 789", "abcd 234 main st 678", "main st", "123 main st"}
	expectedNums := []byte{789 % 256, 678 % 256, 0, 123}
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
