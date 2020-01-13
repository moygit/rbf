package features

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestGetOccurrencePositions(t *testing.T) {
	// given:
	o := OccurrencePositions{true, 3}
	_, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

	// when:
	occurrencePositions := fromStringArray([]string{"", "abcdefgh"})
	blankPositions, abcdefghPositions := occurrencePositions[0], occurrencePositions[1]

	// then:
	if !testSliceIsSingleValue(blankPositions, byte(255)) {
		t.Errorf("got non-255 position for some occurrence position for empty string")
	}

	if !reflect.DeepEqual(abcdefghPositions[:8], []byte{0, 1, 2, 3, 4, 5, 6, 7}) {
		t.Errorf("got incorrect position for some non-infinite occurrence position for string abcdefgh")
	}

	if !testSliceIsSingleValue(abcdefghPositions[8:], byte(255)) {
		t.Errorf("got non-255 position for some occurrence position past h for string abcdefgh")
	}
}

func TestGetLongOccurrencePositions(t *testing.T) {
	// given
	o := OccurrencePositions{true, 3}
	fromString, _ := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

	// when:
	occurrencePositions := fromString(strings.Repeat("a", 400))

	// then:
	if occurrencePositions[0*alphabet_size] != 0 || // first a
		occurrencePositions[1*alphabet_size] != 1 || // second a
		occurrencePositions[2*alphabet_size] != 2 { // third a
		t.Errorf("got invalid occurrence position for letters in string (a**400)+bb")
	}

	if !testSliceIsSingleValue(occurrencePositions[0*alphabet_size+1:1*alphabet_size], byte(255)) ||
		!testSliceIsSingleValue(occurrencePositions[1*alphabet_size+1:2*alphabet_size], byte(255)) ||
		!testSliceIsSingleValue(occurrencePositions[2*alphabet_size+1:3*alphabet_size], byte(255)) {
		t.Errorf("got non-infinite occurrence position for some letter not in string (a**400)+bb")
	}
}

func TestBackwardOccurrencePositions(t *testing.T) {
	// given:
	o := OccurrencePositions{false, 3}
	_, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

	// when:
	occurrencePositions := fromStringArray([]string{"", "abcdefgh"})
	blankPositions, abcdefghPositions := occurrencePositions[0], occurrencePositions[1]

	// then:
	if !testSliceIsSingleValue(blankPositions, byte(255)) {
		t.Errorf("got non-255 position for some backward occurrence position for empty string")
	}

	if !reflect.DeepEqual(abcdefghPositions[:8], []byte{7, 6, 5, 4, 3, 2, 1, 0}) {
		t.Errorf("got incorrect position for some non-infinite backward occurrence position for string abcdefgh")
	}

	if !testSliceIsSingleValue(abcdefghPositions[8:], byte(255)) {
		t.Errorf("got non-255 position for some backward occurrence position past h for string abcdefgh")
	}
}

func TestGetBackwardLongOccurrencePositions(t *testing.T) {
	// given
	o := OccurrencePositions{false, 3}
	fromString, _ := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

	// when:
	occurrencePositions := fromString(strings.Repeat("a", 400) + "bb")

	// then:
	if occurrencePositions[0*alphabet_size+0] != 2 || // first a
		occurrencePositions[1*alphabet_size+0] != 3 || // second a
		occurrencePositions[2*alphabet_size+0] != 4 || // third a
		occurrencePositions[0*alphabet_size+1] != 0 || // first b
		occurrencePositions[1*alphabet_size+1] != 1 || // second b
		occurrencePositions[2*alphabet_size+1] != 255 { // no more b's
		fmt.Printf("%v\n", occurrencePositions)
		t.Errorf("got invalid backward occurrence position for letters in string (a**400)+bb")
	}

	if !testSliceIsSingleValue(occurrencePositions[0*alphabet_size+2:1*alphabet_size], byte(255)) ||
		!testSliceIsSingleValue(occurrencePositions[1*alphabet_size+2:2*alphabet_size], byte(255)) ||
		!testSliceIsSingleValue(occurrencePositions[2*alphabet_size+1:3*alphabet_size], byte(255)) {
		t.Errorf("got non-infinite backward occurrence position for some letter not in string (a**400)+bb")
	}
}
