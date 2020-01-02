package features


import (
    "fmt"
    "reflect"
    "strings"
    "testing"
)


func TestGetOccurrenceCounts(t *testing.T) {
    // given:
    o := OccurrenceCounts{true, 3}
    _, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

    // when:
    occurrenceCounts := fromStringArray([]string{"", "abcdefgh"})
    blankCounts, abcdefghCounts := occurrenceCounts[0], occurrenceCounts[1]

    // then:
    if !testSliceIsSingleValue(blankCounts, byte(255)) {
        t.Errorf("got non-255 count for some occurrence count for empty string")
    }

    if !reflect.DeepEqual(abcdefghCounts[:8], []byte{0, 1, 2, 3, 4, 5, 6, 7}) {
        t.Errorf("got incorrect count for some non-infinite occurrence count for string abcdefgh")
    }

    if !testSliceIsSingleValue(abcdefghCounts[8:], byte(255)) {
        t.Errorf("got non-255 count for some occurrence count past h for string abcdefgh")
    }
}


func TestGetLongOccurrenceCounts(t *testing.T) {
    // given
    o := OccurrenceCounts{true, 3}
    fromString, _ := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

    // when:
    occurrenceCounts := fromString(strings.Repeat("a", 400))

    // then:
    if occurrenceCounts[0 * ALPHABET_SIZE] != 0 ||      // first a
       occurrenceCounts[1 * ALPHABET_SIZE] != 1 ||      // second a
       occurrenceCounts[2 * ALPHABET_SIZE] != 2 {       // third a
        t.Errorf("got invalid occurrence count for letters in string (a**400)+bb")
    }

    if !testSliceIsSingleValue(occurrenceCounts[0 * ALPHABET_SIZE + 1:1 * ALPHABET_SIZE], byte(255)) ||
       !testSliceIsSingleValue(occurrenceCounts[1 * ALPHABET_SIZE + 1:2 * ALPHABET_SIZE], byte(255)) ||
       !testSliceIsSingleValue(occurrenceCounts[2 * ALPHABET_SIZE + 1:3 * ALPHABET_SIZE], byte(255)) {
        t.Errorf("got non-infinite occurrence count for some letter not in string (a**400)+bb")
    }
}


func TestBackwardOccurrenceCounts(t *testing.T) {
    // given:
    o := OccurrenceCounts{false, 3}
    _, fromStringArray := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

    // when:
    occurrenceCounts := fromStringArray([]string{"", "abcdefgh"})
    blankCounts, abcdefghCounts := occurrenceCounts[0], occurrenceCounts[1]

    // then:
    if !testSliceIsSingleValue(blankCounts, byte(255)) {
        t.Errorf("got non-255 count for some backward occurrence count for empty string")
    }

    if !reflect.DeepEqual(abcdefghCounts[:8], []byte{7, 6, 5, 4, 3, 2, 1, 0}) {
        t.Errorf("got incorrect count for some non-infinite backward occurrence count for string abcdefgh")
    }

    if !testSliceIsSingleValue(abcdefghCounts[8:], byte(255)) {
        t.Errorf("got non-255 count for some backward occurrence count past h for string abcdefgh")
    }
}


func TestGetBackwardLongOccurrenceCounts(t *testing.T) {
    // given
    o := OccurrenceCounts{false, 3}
    fromString, _ := MakeFeatureCalculationFunctions([]FeatureSetConfig{o})

    // when:
    occurrenceCounts := fromString(strings.Repeat("a", 400) + "bb")

    // then:
    if occurrenceCounts[0 * ALPHABET_SIZE + 0] != 2 ||      // first a
       occurrenceCounts[1 * ALPHABET_SIZE + 0] != 3 ||      // second a
       occurrenceCounts[2 * ALPHABET_SIZE + 0] != 4 ||      // third a
       occurrenceCounts[0 * ALPHABET_SIZE + 1] != 0 ||      // first b
       occurrenceCounts[1 * ALPHABET_SIZE + 1] != 1 ||      // second b
       occurrenceCounts[2 * ALPHABET_SIZE + 1] != 255 {     // no more b's
        fmt.Printf("%v\n", occurrenceCounts)
        t.Errorf("got invalid backward occurrence count for letters in string (a**400)+bb")
    }

    if !testSliceIsSingleValue(occurrenceCounts[0 * ALPHABET_SIZE + 2:1 * ALPHABET_SIZE], byte(255)) ||
       !testSliceIsSingleValue(occurrenceCounts[1 * ALPHABET_SIZE + 2:2 * ALPHABET_SIZE], byte(255)) ||
       !testSliceIsSingleValue(occurrenceCounts[2 * ALPHABET_SIZE + 1:3 * ALPHABET_SIZE], byte(255)) {
        t.Errorf("got non-infinite backward occurrence count for some letter not in string (a**400)+bb")
    }
}
