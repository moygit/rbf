package features

import (
	"testing"
)

func TestFeatureSetConfig(t *testing.T) {
	// given/when:
	featureSetConfig := []FeatureSetConfig{Followgrams{3}, Followgrams{3}}
	calculateFeatures, _ := MakeFeatureCalculationFunctions(featureSetConfig)
	followgrams := calculateFeatures("abcdefgh")

	// then:
	if !testSliceIsSingleValue(followgrams[1:4], byte(1)) || // a
		!testSliceIsSingleValue(followgrams[39:42], byte(1)) || // b
		!testSliceIsSingleValue(followgrams[77:80], byte(1)) || // c
		!testSliceIsSingleValue(followgrams[115:118], byte(1)) || // d
		!testSliceIsSingleValue(followgrams[153:156], byte(1)) || // e
		!testSliceIsSingleValue(followgrams[191:193], byte(1)) || // f (only 2)
		!testSliceIsSingleValue(followgrams[229:230], byte(1)) || // g (only 1)
		!testSliceIsSingleValue(followgrams[num_followgrams+1:num_followgrams+4], byte(1)) || // a
		!testSliceIsSingleValue(followgrams[num_followgrams+39:num_followgrams+42], byte(1)) || // b
		!testSliceIsSingleValue(followgrams[num_followgrams+77:num_followgrams+80], byte(1)) || // c
		!testSliceIsSingleValue(followgrams[num_followgrams+115:num_followgrams+118], byte(1)) || // d
		!testSliceIsSingleValue(followgrams[num_followgrams+153:num_followgrams+156], byte(1)) || // e
		!testSliceIsSingleValue(followgrams[num_followgrams+191:num_followgrams+193], byte(1)) || // f (only 2)
		!testSliceIsSingleValue(followgrams[num_followgrams+229:num_followgrams+230], byte(1)) || // g (only 1)
		false {
		t.Errorf("abcdefgh 3-followgrams are wrong :-(")
	}
	if !testSliceIsSingleValue(followgrams[0:1], byte(0)) || // before a
		!testSliceIsSingleValue(followgrams[4:39], byte(0)) || // between a and b
		!testSliceIsSingleValue(followgrams[42:77], byte(0)) || // between b and c
		!testSliceIsSingleValue(followgrams[80:115], byte(0)) || // between c and d
		!testSliceIsSingleValue(followgrams[118:153], byte(0)) || // between d and e
		!testSliceIsSingleValue(followgrams[156:191], byte(0)) || // between e and f
		!testSliceIsSingleValue(followgrams[193:229], byte(0)) || // between f and g
		!testSliceIsSingleValue(followgrams[230:num_followgrams], byte(0)) || // after g
		!testSliceIsSingleValue(followgrams[num_followgrams+0:num_followgrams+1], byte(0)) || // before a
		!testSliceIsSingleValue(followgrams[num_followgrams+4:num_followgrams+39], byte(0)) || // between a and b
		!testSliceIsSingleValue(followgrams[num_followgrams+42:num_followgrams+77], byte(0)) || // between b and c
		!testSliceIsSingleValue(followgrams[num_followgrams+80:num_followgrams+115], byte(0)) || // between c and d
		!testSliceIsSingleValue(followgrams[num_followgrams+118:num_followgrams+153], byte(0)) || // between d and e
		!testSliceIsSingleValue(followgrams[num_followgrams+156:num_followgrams+191], byte(0)) || // between e and f
		!testSliceIsSingleValue(followgrams[num_followgrams+193:num_followgrams+229], byte(0)) || // between f and g
		!testSliceIsSingleValue(followgrams[num_followgrams+230:], byte(0)) || // after g
		false {
		t.Errorf("abcdefgh 3-followgrams are wrong :-(")
	}

	// given/when:
	featureSetConfig = []FeatureSetConfig{Followgrams{6}, Followgrams{6}}
	calculateFeatures, _ = MakeFeatureCalculationFunctions(featureSetConfig)
	followgrams = calculateFeatures("aaaaaaaa")

	// then:
	if followgrams[0] != 27 || followgrams[num_followgrams] != 27 {
		t.Errorf("aa count is %d; expected %d", followgrams[0], 27)
	}
	if !testSliceIsSingleValue(followgrams[1:num_followgrams], byte(0)) ||
		!testSliceIsSingleValue(followgrams[num_followgrams+1:], byte(0)) {
		t.Errorf("got non-zero count for some non-aa value")
	}
}
