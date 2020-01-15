package features

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetOccurrenceCounts(t *testing.T) {
	// given:
	o := OccurrenceCounts{2}
	_, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

	// when:
	occurrenceCounts := fromStringArray([]string{"", "aaabbcddd"})
	blankCounts, aaabbcdddCounts := occurrenceCounts[0], occurrenceCounts[1]

	// then:
	if !testSliceIsSingleValue(blankCounts, byte(0)) {
		t.Errorf("got non-255 count for some occurrence count for empty string")
	}

	if !reflect.DeepEqual(aaabbcdddCounts[:4], []byte{3, 2, 1, 3}) ||
		!reflect.DeepEqual(aaabbcdddCounts[alphabet_size:alphabet_size+4], []byte{3, 2, 1, 3}) {
		t.Errorf("got incorrect count for some non-zero occurrence count for string aaabbcddd")
	}

	if !testSliceIsSingleValue(aaabbcdddCounts[4:alphabet_size], byte(0)) ||
		!testSliceIsSingleValue(aaabbcdddCounts[alphabet_size+4:], byte(0)) {
		t.Errorf("got non-255 count for some occurrence count past d string aaabbcddd")
	}
}

func TestGetLongOccurrenceCounts(t *testing.T) {
	// given
	o := OccurrenceCounts{1}
	fromString, _ := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

	// when:
	occurrenceCounts := fromString(strings.Repeat("a", 400))

	// then:
	if occurrenceCounts[0] != 255 {
		t.Errorf("got invalid occurrence count for letters in string (a**400)")
	}
}
