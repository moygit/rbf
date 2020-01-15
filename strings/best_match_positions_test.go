package strings

import (
    "reflect"
    "testing"
)

func TestMin(t *testing.T) {
	// given:
    L := []int{5, 8, 1, 6, 31, 14, 19}
	// when:
    minVal := min(L...)
	// then:
	if minVal != 1 {
		t.Errorf("min value of %v should not be %d", L, minVal)
	}
}

func TestMax(t *testing.T) {
	// given:
    L := []int{5, 8, 1, 6, 31, 14, 19}
	// when:
    maxVal := max(L...)
	// then:
	if maxVal != 31 {
		t.Errorf("max value of %v should not be %d", L, maxVal)
	}
}

func TestGetCounts(t *testing.T) {
	// given:
    s := "this is a test"
	// when:
    counts := getCounts([]byte(s))
	// then:
	if counts[1-1] != 1   ||    // a
	   counts[5-1] != 1   ||    // e
	   counts[8-1] != 1   ||    // h
	   counts[9-1] != 2   ||    // i
	   counts[19-1] != 3  ||    // s
	   counts[20-1] != 3  ||    // t
	   counts[37-1] != 3 {      // space
		t.Errorf("incorrect counts for string 'this is a test'\n")
	}
}

func TestGetIntersection(t *testing.T) {
    compare := func(s1, s2 string, expected int) {
        // when:
        intersection := getIntersection(getCounts([]byte(s1)), getCounts([]byte(s2)))
        // then:
        if intersection != expected {
            t.Errorf("incorrect intersection for strings '%s' and '%s'\n", s1, s2)
        }
    }

	// given:
    ref := "this is a test"
    s1 := ""
    s2 := "this is a test"
    s3 := "also this is a"
    s4 := "also this is a test too"
    s5 := "this too is a "

    compare(ref, s1, 0)
    compare(ref, s2, 14)
    compare(ref, s3, 11)
    compare(ref, s4, 14)
    compare(ref, s5, 11)
}

func TestGetLocalMaximaAboveThreshold(t *testing.T) {
    // given:
    L := []int{17, 0, 1, 0, 5, 7, 9, 9, 7, 10, 10, 11, 12}
    threshold := 2
    // when:
    maxima := getLocalMaximaAboveThreshold(L, threshold)
    // then:
    // Expected maxima: 0 (17), 6 (9), 7 (9), 12 (12)
    // BUT ALSO: 9 (10) because we don't want to waste cycles looking that far ahead
    expectedMaxima := []int{0, 6, 7, 9, 12}
	if !reflect.DeepEqual(maxima, expectedMaxima) {
        t.Errorf("Got unexpected local maxima (actually argMax): %v\n", maxima)
	}
}

func TestGetBestMatchPositions(t *testing.T) {
    compare := func(s1, s2 string, expectedPositions []int) {
        positions := GetBestMatchPositions(s1, s2)
        if !reflect.DeepEqual(positions, expectedPositions) {
            t.Errorf("Got unexpected max intersection positions: %v, expected: %v\n", positions, expectedPositions)
        }
    }

    // test cases:
    ref := "abcd"
    s1 := ""
    s2 := "a"
    s3 := "ab"
    s4 := "abcde"
    s5 := "abcdefgh"
    s6 := "abcdefgh abcd"

    compare(ref, s1, []int{})
    compare(ref, s2, []int{})
    compare(ref, s3, []int{})
    compare(ref, s4, []int{0})
    compare(ref, s5, []int{0})
    compare(ref, s6, []int{0, 9})
}
