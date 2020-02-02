package features

import (
	"testing"
)

func TestFeatureSetConfig(t *testing.T) {
	// given/when:
	featureSetConfigStr := `
- feature_type: followgrams
  window_size: 3
- feature_type: followgrams
  window_size: 3
`
	calculateFeatures, _ := GetFeatureCalcFuncs(featureSetConfigStr)
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
	featureSetConfigStr = `
- feature_type: followgrams
  window_size: 6
- feature_type: followgrams
  window_size: 6
`
	calculateFeatures, _ = GetFeatureCalcFuncs(featureSetConfigStr)
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

func catchPanicOrElse(t *testing.T, msg string) {
	if r := recover(); r == nil {
		t.Errorf(msg)
	}
}

func TestReadSerializedConfigFailsOnBadInput(t *testing.T) {
	featureConfig := "- key: n1\n  val1: v1\n  val2: v2\n- key: n2\n  val3: 3\n  val4: 4\n"
	defer catchPanicOrElse(t, "Config without key 'feature_type' should have panicked but didn't.")
	getConfigsFromYaml(featureConfig)
}

func TestReadSerializedConfigWithUnknownType(t *testing.T) {
	featureConfig := "- feature_type: unknown_type_should_cause_panic\n  val1: v1\n  val2: v2\n"
	defer catchPanicOrElse(t, "Unknown feature type unknown_type_should_cause_panic did not panic")
	getConfigsFromYaml(featureConfig)
}

func TestReadSerializedConfigFromString(t *testing.T) {
	// given
	expectedCount := 6
	featureConfig := `
- feature_type: first_number
  count: "17"   # w/ quotes here, w/o below; both should work
- feature_type: last_number
  # just use default count
- feature_type: followgrams
  window_size: 3
- feature_type: bigrams
  allow_repeats: true
- feature_type: occurrence_positions
  direction_is_head: true
  num_occurrences: 5
- feature_type: followgrams
`
	// when
	configs := getConfigsFromYaml(featureConfig)
	// then
	if count := len(configs); count != expectedCount {
		t.Errorf("deserialized featureSetConfigs slice [%v] should have contained %d values but instead contains %d", configs, expectedCount, count)
	}
	if f, ok := configs[0].(firstNumber); !ok || f.Count != 17 {
		t.Errorf("expected configs[0] (%v) to be firstNumber{17}", configs[0])
	}
	if f, ok := configs[1].(lastNumber); !ok || f.Count != last_number_default_count {
		t.Errorf("expected configs[1] (%v) to be lastNumber{%d}", configs[1], last_number_default_count)
	}
	if f, ok := configs[2].(followgrams); !ok || f.WindowSize != 3 {
		t.Errorf("expected configs[2] (%v) to be followgrams{3}", configs[2])
	}
	if f, ok := configs[3].(bigrams); !ok || f.maxBigramCount != 255 {
		t.Errorf("expected configs[3] (%v) to be bigrams{255}", configs[3])
	}
	if f, ok := configs[4].(occurrencePositions); !ok || !f.DirectionIsHead || f.NumberOfOccurrences != 5 {
		t.Errorf("expected configs[4] (%v) to be occurrencePositions{true, 5}", configs[4])
	}
	if f, ok := configs[5].(followgrams); !ok || f.WindowSize != followgram_default_window_size {
		t.Errorf("expected configs[5] (%v) to be followgrams{%d}", configs[5], followgram_default_window_size)
	}
}
