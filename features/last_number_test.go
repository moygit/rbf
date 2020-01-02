package features


import "testing"


func TestGetLastNumber(t *testing.T) {
    // given:
    strs := []string{"123 main st 789", "abcd 234 main st 678", "main st", "123 main st"}
    expectedNums := []byte{789 % 256, 678 % 256, 0, 123}
    // when/then:
    for i := 0; i < 4; i++ {
        n := GetLastNumber(strs[i])
        if n != expectedNums[i] {
            t.Errorf("bad last-num %d for string %s, expected %d\n", n, strs[i], expectedNums[i])
        }
    }
}


func TestGetLastNumberFeature(t *testing.T) {
    // given:
    fn := DefaultLastNumber
    _, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{fn})
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

